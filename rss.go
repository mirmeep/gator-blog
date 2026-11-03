package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/xml"
	"context"
	"html"
	"time"
	"github.com/mirmeep/gator-blog/internal/database"
	"github.com/google/uuid"
	"database/sql"
	"github.com/lib/pq"
	"strconv"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer res.Body.Close()
	if res.StatusCode > 299 {
		return &RSSFeed{}, fmt.Errorf("RSS feed not found")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rss RSSFeed
	if err := xml.Unmarshal(data, &rss); err != nil {
		return &RSSFeed{}, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	
	for index := range rss.Channel.Item {
		rss.Channel.Item[index].Title = html.UnescapeString(rss.Channel.Item[index].Title)
		rss.Channel.Item[index].Description = html.UnescapeString(rss.Channel.Item[index].Description)
	}

	return &rss, nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("command requires one arg, time duration")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	fmt.Println("Time between reqs is:", timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("could not parse duration: %w", err)
	}

	if timeBetweenReqs <= 0 {
		return fmt.Errorf("not a valid time duration")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("command requires two params, name and url")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return fmt.Errorf("feed already exists")
	}

	if err != nil {
		return fmt.Errorf("cannot create feed: %w", err)
	}
	fmt.Println("Added feed for", user.Name)
	fmt.Printf("* ID:            	%s\n", feed.ID)
	fmt.Printf("* Created:       	%v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       	%v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          	%s\n", feed.Name)
	fmt.Printf("* URL:           	%s\n", feed.Url)
	fmt.Printf("* User ID:          %s\n", feed.UserID)
	fmt.Println("=====================================")


	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("cannot insert into table: %w", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("command feeds doesn't take params")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Errorf("cannot fetch feeds: %w", err)
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(),feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		fmt.Printf("* ID:            %s\n", feed.ID)
		fmt.Printf("* Created:       %v\n", feed.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* User:          %s\n", user.Name)
		fmt.Println("=====================================")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("command follow requires one arg, url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("cannot get user feeds: %w", err)
	}

	followFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("cannot insert into table: %w", err)
	}

	fmt.Printf("* Name:          %s\n", followFeed.FeedName)
	fmt.Printf("* User:          %s\n", followFeed.UserName)
	fmt.Println("=====================================")

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("command following doesn't take args")
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("cannot get user follows: %w", err)	
	}
	fmt.Printf("%s is following:\n", user.Name)
	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}
	fmt.Println("=====================================")
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("command unfollow requires one arg, url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	unfollowed, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams {
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}

	fmt.Println("Successfully unfollowed feed:")
	fmt.Printf("* ID:            %s\n", unfollowed.ID)
	fmt.Printf("* Created:       %v\n", unfollowed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", unfollowed.UpdatedAt)
	fmt.Printf("* UserID:        %s\n", unfollowed.UserID)
	fmt.Printf("* FeedID:        %s\n", unfollowed.FeedID)
	fmt.Println("=====================================")

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.conf.Username)
		if err != nil {
			return fmt.Errorf("cannot find current user: %w", err)
		}	
		return handler(s, cmd, currentUser)
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("cannot fetch next feed: %w", err) 
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("could not update feed last fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	fmt.Println("-------------------------------------")	
	fmt.Println("Fetched feed for", feed.Channel.Title)
	fmt.Println("-------------------------------------")
	for index := range feed.Channel.Item {
		var parsedTime time.Time
		var err error
		parsedTime, err = time.Parse(time.RFC1123Z, feed.Channel.Item[index].PubDate)
		if err != nil {
			parsedTime, err = time.Parse(time.RFC1123, feed.Channel.Item[index].PubDate)	
			
		}
		if err != nil {
			parsedTime, err = time.Parse(time.RFC3339, feed.Channel.Item[index].PubDate)
			
		}
		publishedDate := sql.NullTime {
			Time: parsedTime,
			Valid: err == nil,
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: feed.Channel.Item[index].Title,
			Url: feed.Channel.Item[index].Link,
			Description: sql.NullString{
				String: feed.Channel.Item[index].Description,
				Valid: true,
			},
			PublishedAt: publishedDate,
			FeedID: nextFeed.ID,
		})

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Errorf("post already exists in database: %w", err)
			continue
		}
		fmt.Errorf("couldn't create post: %w", err)
		continue
	}
	fmt.Println("=====================================")

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("could not parse string to int: %w", err)
		}
	} 

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID, 
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}
	
	fmt.Println("=====================================")	
	fmt.Println("Printing posts for", user.Name)
	fmt.Println("=====================================")	
	for _, post := range posts {
		fmt.Printf("* ID:            		%s\n", post.ID)
		fmt.Printf("* Created:       		%v\n", post.CreatedAt)
		fmt.Printf("* Updated:       		%v\n", post.UpdatedAt)
		fmt.Printf("* Title:          		%s\n", post.Title)
		fmt.Printf("* URL:           		%s\n", post.Url)
		fmt.Printf("* Description:          %s\n", post.Description.String)
		fmt.Printf("* Published:			%s\n", post.PublishedAt.Time)
		fmt.Printf("* Feed ID:				%s\n", post.FeedID)		
		fmt.Println("-------------------------------------")
	}

	return nil
}
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ArchieSpinos/grpc-go/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "Archie",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Cloudn't create blog: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()
	fmt.Printf("This is the blog: %v", blog)

	// read blog
	fmt.Println("reading the blog")
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "thistest"})
	if err != nil {
		fmt.Printf("error happened while reading: %v", err2)
	}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("error happened while reading: %v", readBlogErr)
	}

	fmt.Printf("Blow was read: %v", readBlogRes)

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was read: %v", updateRes)

	// deleteBlogRes, deleteBlogErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{Blog: blog})

	// if deleteBlogErr != nil {
	// 	fmt.Printf("Error happened while deleting: %v \n", deleteBlogErr)
	// }
	// fmt.Printf("Blog was deleted: %v", deleteBlogRes)

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}

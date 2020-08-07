  1 package main
  2 
  3 import (
  4 >       "context"
  5 >       "fmt"
  6 
  7 >       vision "cloud.google.com/go/vision/apiv1"
  8 )
  9 
 10 func annotate(uri string) (float32, error) {
 11 >       ctx := context.Background()
 12 
 13 >       client, err := vision.NewImageAnnotatorClient(ctx)
 14 >       if err != nil {
 15 >       >       return 0.0, err
 16 >       }
 17 
 18 >       image := vision.NewImageFromURI(uri)
 19 
 20 >       annotations, err := client.DetectFaces(ctx, image, nil, 1)
 21 >       if err != nil {
 22 >       >       return 0.0, err
 23 >       }
 24 
 25 >       //if no faces in the image
 26 >       if len(annotations) == 0 {
 27 >       >       fmt.Println("No face detected")
 28 >       >       return 0.0, nil
 29 >       }
 30 
 31 >       return annotations[0].DetectionConfidence, nil
 32 }
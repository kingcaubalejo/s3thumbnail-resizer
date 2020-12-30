# s3ThumbnailResizer

Out of boredom and procastination. I created a thumbnail resizer, which resize the image based on the width and length of your choice. It also gives you
an option to upload it to s3 and dynamo db for the metadata of the image.

## Getting Started

I basically used cloud for this mini-project,  using Lanczos3 https://en.wikipedia.org/wiki/Lanczos_resampling resizer


### Prerequisites

Library used

* https://golang.org/doc/install
* https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html
* https://github.com/nfnt/resize

## Deployment

cd to your project folder<br/>
go build -o name <br/>
./name <br/> 

## Built With

* Love, nah just kidding. Built with AWS Services and Go lang.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

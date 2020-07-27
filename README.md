# Bucket-Details

![Go](https://github.com/henriquegibin/bucket-details/workflows/Go/badge.svg?branch=master)

- [What is](#What-is)
- [Requirements](#Requirements)
- [Build And Tests](#Build-And-Test)
  - [Build](#Build)
  - [Test](#Test)
- [Usage](#Usage)

## What is

Bucket-Details is a tool to get some information from your S3 buckets.

## Requirements

Bucket-Details is a golang project, so ou just need to download the binarie and you are ready to go.
But if you want to compile by yourself you can use:

- Docker
- Go (1.13+)

And you will need to export your AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY as environment variables.

## Build And Test

### Build

To build this binarie you have some options. Using Go:

```bash
go mod download # Download all project dependencies
go build -o bucket-details # Build the binarie and let inside this folder
```

Using docker:

```bash
docker build . -t bucket-details
docker run -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY bucket-details # Run without any flag
docker run -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY harry ./bucket-details --life-cycle # Run with flags
```

Using docker-compose:

```bash
docker-compose run app # Run without any flag
docker-compose run app ./bucket-details --life-cycle um # Run with flags
```

### Test

To run tests you can use go and docker. Using Go:

```bash
go test -cover ./src/* # This will output the coverage and the test results
```

To use docker, just run the build step. Tests will run automatically during the process.

## Usage

After downloading the binarie (or compiling), you can execute doing this:

```bash
./bucket-details
```

After this, all your buckets will be scanned and the basic info will appear in your stdout.
The basic info is:

- Bucket Name
- Creation Date
- Number of files
- Total size of files in KB
- Last modified date of the most recent file
- And how much does it cost

It's important to know, the cost will be $0 if your buckets don't have a tag `Name`.
I use the cost Explorer to retrieve the exact cost. But to filter, this tag is needed.

You can also pass a set of flags to get additional info and filter buckets.
Here is a list with all flags(you can pass `--help` to list in your terminal):

| Flag                    | Description                                                          |
| ----------------------- | -------------------------------------------------------------------- |
| life-cycle value        |  Pass this flag to retreive the bucket life cycle (default: "false") |
| bucket-acl value        |  Pass this flag to retreive the bucket bucket acl (default: "false") |
| bucket-Encryption value |  Pass this flag to retreive the bucket encryption (default: "false") |
| bucket-location value   |  Pass this flag to retreive the bucket location (default: "false")   |
| bucket-tagging value    |  Pass this flag to retreive the bucket tagging (default: "false")    |

---

### Filter

To filter by bucket name you have 3 options

- prefix
- suffix
- regex

To use filters you need to pass the flag that indicate `type` and the flag with `value`:

```bash
./bucket-details --filter-type prefix --filter-value ab
./bucket-details --filter-type suffix --filter-value .br
./bucket-details --filter-type regex --filter-value \^\[a-z\]\+\.\[a-z\]\+\-\[a-z\]\+\.com\.br
```

## Improvements

I'm not happy with the way the code is organized. It's look a mess. And the performance could be better using `go routines` to retrieve all metadatas.

If you run inside a private ec2 instance, its faster because of the aws internal communication structure, but i'm not happy with the results yet.
I will improving using this project to study Go and learn how to use go routines.

Another thing that could improve the project is to use it together with a chatbot.
The chatbot could create a friendly interface between the user and the app.

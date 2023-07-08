# awsutils/s3

The `s3` package is a collection of utility functions
designed to simplify common s3 tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CreateBucket(*s3.S3, string)

```go
CreateBucket(*s3.S3, string) error
```

CreateBucket creates a bucket with the input
bucketName.

Parameters:
client: An AWS S3 client.
bucketName: The name of the bucket to create.

Returns:
error: An error if the bucket could not be created.

---

### CreateConnection()

```go
CreateConnection() Connection
```

CreateConnection creates a connection
with S3 and returns it.

Returns:
An s3.Connection struct containing the AWS S3 client and AWS session.

---

### DestroyBucket(*s3.S3, string)

```go
DestroyBucket(*s3.S3, string) error
```

DestroyBucket destroys a bucket with the input
bucketName.

Parameters:
client: An AWS S3 client.
bucketName: The name of the bucket to destroy.

Returns:
error: An error if the bucket could not be destroyed.

---

### DownloadBucketFile(*session.Session, string, string, string)

```go
DownloadBucketFile(*session.Session, string, string, string) string, error
```

DownloadBucketFile downloads a file found at the object key specified with
the input objectKey from the bucket specified by bucketName, and writes it to downloadFP.

Parameters:
sess: An AWS session.
bucketName: The name of the bucket to download from.
objectKey: The key of the object to download.
downloadFP: The file path to write the downloaded file to.

Returns:
string: The name of the downloaded file.
error: An error if the file could not be downloaded.

---

### EmptyBucket(*s3.S3, string)

```go
EmptyBucket(*s3.S3, string) error
```

EmptyBucket deletes everything found in the input bucketName.

Parameters:
client: An AWS S3 client.
bucketName: The name of the bucket to empty.

Returns:
error: An error if the bucket could not be emptied.

---

### GetBuckets(*s3.S3)

```go
GetBuckets(*s3.S3) []*s3.Bucket, error
```

GetBuckets returns all s3 buckets
that the input client has access to.

Parameters:
client: An AWS S3 client.

Returns:
[]*s3.Bucket: A slice of pointers to AWS S3 bucket structs that the client has access to.
error: An error if the buckets could not be listed.

---

### UploadBucketDir(*session.Session, string, string)

```go
UploadBucketDir(*session.Session, string, string) error
```

UploadBucketDir uploads a directory specified by dirPath
to the bucket specified by bucketName.

Parameters:
sess: An AWS session.
bucketName: The name of the bucket to upload to.
dirPath: The file path of the directory to upload.

Returns:
error: An error if the directory could not be uploaded.

---

### UploadBucketFile(*session.Session, string, string)

```go
UploadBucketFile(*session.Session, string, string) error
```

UploadBucketFile uploads a file found at the file path specified with (uploadFP)
to the input bucketName.

Parameters:
sess: An AWS session.
bucketName: The name of the bucket to upload to.
uploadFP: The file path of the file to upload.

Returns:
error: An error if the file could not be uploaded.

---

## Installation

To use the awsutils/s3 package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/s3
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/s3"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `s3/awsutils`:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes,
please open an issue first to discuss what
you would like to change.

---

## License

This project is licensed under the MIT
License - see the [LICENSE](../LICENSE)
file for details.

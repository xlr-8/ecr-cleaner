# ecr-cleaner
Deletes old images from ecr

This clean up a specific repository as well as all repos within an aws account.


### Algorithm
1. Retrieve repo from ecr
2. Get repo images
3. Add all images without tags to deletion
4. Sort the remaining images in alphanumeric order
5. Add n oldest images to deletion
6. Delete images from the repository


### Default values
`aws.region = eu-central-1`
`dry-run = false`
`amount-to-keep = 100`


### Examples
clean up all repos

`ecr-cleaner -aws.region=eu-west-1`

clean up my-awesome-repo

`ecr-cleaner -aws.region=eu-west-1 -repository my-awesome-repo`

go for a dry run

`ecr-cleaner -aws.region=eu-west-1 -repository my-awesome-repo -dry-run true`

leave n images in repo

`ecr-cleaner -aws.region=eu-west-1 -repository my-awesome-repo -amount-to-keep 5`

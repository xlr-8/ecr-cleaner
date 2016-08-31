variable "aws_region" {
  default = "eu-central-1"
}

variable "cron" {
  default = "0 3 1 * ? *"
}

variable "repository" {
  default = ""
}

variable "amount_of_images_to_keep" {
  default = ""
}

variable "dry_run" {
  default = "true"
}

variable "repo_region" {
  default = "eu-west-1"
}

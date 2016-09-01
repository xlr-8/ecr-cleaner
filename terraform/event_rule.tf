resource "aws_cloudwatch_event_rule" "scheduled_ecr_cleaner" {
  name                = "scheduled_ecr_cleaner"
  description         = "Run ecr cleaner"
  schedule_expression = "cron(${var.cron})"
}

resource "aws_cloudwatch_event_target" "scheduled_ecr_cleaner_lambda" {
  rule  = "${aws_cloudwatch_event_rule.scheduled_ecr_cleaner.name}"
  arn   = "${aws_lambda_function.lambda.arn}"
  input = "${data.template_file.event_json.rendered}"
}

data "template_file" "event_json" {
  template = "${file("${path.module}/event_json.json")}"

  vars = {
    aws_region = "${var.repo_region}"
    dry_run    = "${var.dry_run}"
    repository = "${var.repository}"
  }
}

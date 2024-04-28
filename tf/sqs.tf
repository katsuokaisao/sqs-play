resource "aws_sqs_queue" "example" {
  name_prefix                = "example"
  visibility_timeout_seconds = "30"   # 30 seconds
  message_retention_seconds  = 345600 # 4 days
  max_message_size           = 262144 # 256 KiB (max)
  delay_seconds              = 0
  receive_wait_time_seconds  = 20 # 20 seconds (max)

  #   TODO
  #   redrive_policy = jsonencode({
  #     deadLetterTargetArn = aws_sqs_queue.terraform_queue_deadletter.arn
  #     maxReceiveCount     = 4
  #   })
  #   redrive_policy

  fifo_queue                  = true
  content_based_deduplication = true
}

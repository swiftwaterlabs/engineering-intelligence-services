resource "aws_sqs_queue" "signal_ingestion" {
  name                  = "${local.service_name}_signal_ingestion"
  fifo_queue            = false
}

resource "aws_sqs_queue" "webhook_event" {
  name                  = "${local.service_name}_webhook_event"
  fifo_queue            = false
}
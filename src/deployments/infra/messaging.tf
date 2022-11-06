resource "aws_sqs_queue" "signal_ingestion" {
  name                  = "${local.service_name}_signal_ingestion"
  fifo_queue            = false
}
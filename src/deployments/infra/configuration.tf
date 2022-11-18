resource "aws_secretsmanager_secret" "hosts_table" {
  name = "${local.service_name}_directories_table"
}

resource "aws_secretsmanager_secret_version" "hosts_table" {
  secret_id     = aws_secretsmanager_secret.hosts_table.id
  secret_string = aws_dynamodb_table.hosts.name
}

resource "aws_secretsmanager_secret" "event_sources_table" {
  name = "${local.service_name}_event_sources_table"
}

resource "aws_secretsmanager_secret_version" "event_sources_table" {
  secret_id     = aws_secretsmanager_secret.event_sources_table.id
  secret_string = aws_dynamodb_table.event_sources.name
}

resource "aws_secretsmanager_secret" "ingestion_queue" {
  name = "${local.service_name}_ingestion_queue"
}

resource "aws_secretsmanager_secret_version" "ingestion_queue" {
  secret_id     = aws_secretsmanager_secret.ingestion_queue.id
  secret_string = aws_sqs_queue.signal_ingestion.name
}

resource "aws_secretsmanager_secret" "webhook_event_queue" {
  name = "${local.service_name}_webhook_event_queue"
}

resource "aws_secretsmanager_secret_version" "webhook_event_queue" {
  secret_id     = aws_secretsmanager_secret.webhook_event_queue.id
  secret_string = aws_sqs_queue.webhook_event.name
}

resource "aws_secretsmanager_secret" "blob_store" {
  name = "${local.service_name}_blob_store"
}

resource "aws_secretsmanager_secret_version" "blob_store" {
  secret_id     = aws_secretsmanager_secret.blob_store.id
  secret_string = local.signal_bucket_name
}

resource "aws_secretsmanager_secret" "aws_region" {
  name = "${local.service_name}_aws_region"
}

resource "aws_secretsmanager_secret_version" "aws_region" {
  secret_id     = aws_secretsmanager_secret.aws_region.id
  secret_string = var.aws_region
}
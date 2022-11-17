data "archive_file" "webhookreceiver_lambda_zip" {
  type        = "zip"
  source_file = "../../cmd/lambda-webhookreceiver/main"
  output_path = "webhookreceiver_main.zip"
}

resource "aws_lambda_function" "webhook_receiver" {
  function_name = "${local.service_name}_webhook_receiver"

  role = aws_iam_role.lambda_exec.arn

  filename          = data.archive_file.webhookreceiver_lambda_zip.output_path
  handler           = "main"
  source_code_hash  = filebase64sha256(data.archive_file.webhookreceiver_lambda_zip.output_path)
  runtime           = "go1.x"

  environment {
    variables = {
      aws_region = var.aws_region
    }
  }
  
}

resource "aws_cloudwatch_log_group" "webhook_receiver" {
  name = "/aws/lambda/${aws_lambda_function.webhook_receiver.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "webhook_receiver" {
  api_id = aws_apigatewayv2_api.lambda_gateway.id

  integration_uri    = aws_lambda_function.webhook_receiver.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  
}

resource "aws_apigatewayv2_route" "webhook_receiver" {
  api_id = aws_apigatewayv2_api.lambda_gateway.id

  route_key = "POST /signal"
  target    = "integrations/${aws_apigatewayv2_integration.webhook_receiver.id}"
  
}

resource "aws_lambda_permission" "webhook_receiver" {

  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.webhook_receiver.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda_gateway.execution_arn}/*/*"
}
resource "aws_glue_catalog_database" "object_db" {
  name = local.service_name
}

resource "aws_glue_catalog_table" "repository" {
  name          = "repository"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/repository"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }


    columns {
      name = "organization"
      type = "struct<id:string,name:string,host:string,hosttype:string,url:string>"
    }

    columns {
      name = "name"
      type = "string"
    }

    columns {
      name = "url"
      type = "string"
    }

    columns {
      name = "defaultbranch"
      type = "string"
    }

    columns {
      name = "visibility"
      type = "string"
    }

    columns {
      name = "createdat"
      type = "string"
    }

    columns {
      name = "updatedat"
      type = "string"
    }
  }
}

resource "aws_glue_catalog_table" "repository_owner" {
  name          = "repository_owner"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/repository-owner"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }

    columns {
      name = "repository"
      type = "struct<id:string,name:string,organization:struct<id:string,name:string,host:string,hosttype:string,url:string>>"
    }

    columns {
      name = "pattern"
      type = "string"
    }

    columns {
      name = "owner"
      type = "string"
    }

    columns {
      name = "parentowner"
      type = "string"
    }

  }
}

resource "aws_glue_catalog_table" "pull_request" {
  name          = "pull_request"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/pull-request"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }


    columns {
      name = "repository"
      type = "struct<id:string,name:string,organization:struct<id:string,name:string,host:string,hosttype:string,url:string>>"
    }

    columns {
      name = "targetbranch"
      type = "string"
    }

    columns {
      name = "url"
      type = "string"
    }

    columns {
      name = "title"
      type = "string"
    }

    columns {
      name = "createdby"
      type = "string"
    }

    columns {
      name = "createdat"
      type = "string"
    }

    columns {
      name = "closedat"
      type = "string"
    }

    columns {
      name = "ismerged"
      type = "boolean"
    }

    columns {
      name = "status"
      type = "string"
    }

    columns {
      name="reviews"
      type = "array<struct<reviewer:string,status:string,reviewedat:string>>"
    }

    columns {
      name="files"
      type="array<string>"
    }

  }
}

resource "aws_glue_catalog_table" "branch_rule" {
  name          = "branch_rule"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/branch-rule"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }


    columns {
      name = "repository"
      type = "struct<id:string,name:string,organization:struct<id:string,name:string,host:string,hosttype:string,url:string>>"
    }

    columns {
      name = "branch"
      type = "string"
    }

    columns {
      name = "allowforcepush"
      type = "boolean"
    }

    columns {
      name = "requirepullrequest"
      type = "boolean"
    }

    columns {
      name = "requirepullrequestapprovals"
      type = "boolean"
    }

    columns {
      name = "requiredpullrequestapprovers"
      type = "int"
    }

    columns {
      name = "includeadministrators"
      type = "boolean"
    }

  }
}

resource "aws_glue_catalog_table" "webhook" {
  name          = "webhook"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/webhook"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }

    columns {
      name = "name"
      type = "string"
    }

    columns {
      name = "source"
      type = "string"
    }

    columns{
      name = "organization"
      type = "struct<id:string,name:string,host:string,hosttype:string,url:string>"
    }

    columns {
      name = "repository"
      type = "struct<id:string,name:string,organization:struct<id:string,name:string,host:string,hosttype:string,url:string>>"
    }

    columns {
      name = "target"
      type = "string"
    }

    columns {
      name = "events"
      type = "array<string>"
    }

    columns {
      name = "active"
      type = "boolean"
    }

  }
}

resource "aws_glue_catalog_table" "test_result" {
  name          = "test_result"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/testresult"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }

    columns {
      name = "project"
      type = "string"
    }

    columns {
      name = "analyzedat"
      type = "string"
    }

    columns{
      name = "metrics"
      type = "struct<coverage:string,lines_to_cover:string,uncovered_lines:string>"
    }


  }
}
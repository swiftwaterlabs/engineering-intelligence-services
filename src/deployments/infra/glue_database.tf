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
      serialization_library = "org.apache.hive.hcatalog.data.JsonSerDe"

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
      serialization_library = "org.apache.hive.hcatalog.data.JsonSerDe"

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
      name = "repositoryname"
      type = "string"
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
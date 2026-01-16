env "local" {
  src = "file://infrastructure/pgsql/migrations"
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://infrastructure/pgsql/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://infrastructure/pgsql/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

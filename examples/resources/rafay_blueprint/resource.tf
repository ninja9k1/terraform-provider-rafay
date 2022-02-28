resource "rafay_blueprint" "tfdemoblueprint1" {
  metadata {
    name    = "tfdemoblueprint1"
    project = "upgrade"
    labels = {
      env  = "dev"
      name = "app"
    }
  }
  spec {
    version = "1.1"
    base {
      name    = "default"
      version = "1.11.0"
    }
    default_addons {
      enable_ingress    = true
      enable_logging    = true
      enable_monitoring = false
      enable_vm         = false
      monitoring {
        metrics_server {
          enabled = false
          discovery {}
        }
        helm_exporter {
          enabled = false
        }

        kube_state_metrics {
          enabled = false
        }

        node_exporter {
          enabled = false
        }

        prometheus_adapter {
          enabled = false
        }

        resources {
        }
      }
    }

    drift {
      action  = "Deny"
      enabled = true
    }

    sharing {
      enabled = false
      projects {
        name = "demo"
      }
    }
  }
}


resource "rafay_blueprint" "tfdemoblueprint2" {
  metadata {
    annotations = {}
    labels = {
      "env"  = "dev"
      "name" = "app"
    }
    name    = "tfdemoblueprint2"
    project = "upgrade"
  }

  spec {
    version = "1.1"

    base {
      name    = "default"
      version = "1.11.0"
    }

    custom_addons {
      depends_on = []
      name       = "tomcat1"
      version    = "v1"
    }

    custom_addons {
      depends_on = []
      name       = "gold-pinger"
      version    = "v0"
    }

    default_addons {
      enable_ingress    = true
      enable_logging    = true
      enable_monitoring = true
      enable_vm         = false

      monitoring {
        helm_exporter {
          enabled = false
        }

        kube_state_metrics {
          enabled = false
        }

        metrics_server {
          enabled = false

          discovery {}
        }

        node_exporter {
          enabled = false
        }

        prometheus_adapter {
          enabled = false
        }

        resources {
        }
      }
    }

    drift {
      action  = "Deny"
      enabled = true
    }

    sharing {
      enabled = false
    }
  }
}

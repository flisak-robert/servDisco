Filed based service discovery that works with Prometheus

Persistence will probably be needed so in order to achieve that, create these to files first:

`touch targets.json targets-ssl.json`

And run the container with the following command:

`docker run -itd -p 4000:4000 -v $(pwd)/targets.json:/app/targets.json -v $(pwd)/targets-ssl.json:/app/targets-ssl.json servdisco`

The service inside the container is running on port 4000. Remember to correctly mount the same files to prometheus container and update the prometheus.yml config file.

To add a host using a POST request, use the `/add` endpoint. Include the following JSON body:

```
{
    "targets": ["{{ target_IP/target_hostname }}:{{ target_port }}"],
    "labels": {
        "job": "{{ monitoring_job_name }}",
        "env": "{{ project_environment }}",
        "__metrics_path__": "{{ monitoring_metrics_path }}"
    }
}
```

Example of using curl to add a new target:

```curl --header "Content-Type: application/json" -X POST -d '{"targets":["localhost:8000"], "labels": { "job": "cadvisor", "env":"dev", "__metrics_path__": "/metrics"}}' localhost:4000/add/targets```

Works with blackbox-exporter for SSL certificate expiration date monitoring.
To add targets to SSL monitoring, use `/add/targets-ssl` endpoint.
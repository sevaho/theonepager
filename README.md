# theonepager

A homepage for your virtual home in one page, or more.

# Quickstart

## Docker

Write a `config.yaml` 

```yaml
---
applications:
  - name: Grafana
    category: Monitoring
    link: https://grafana.internal
    description: |
      This
      is
      more
      info
```


```bash
docker run -p3000:3000 -v $(pwd)/config.yaml:/config.yaml sevaho/theonepager --serve -c "/config.yaml"
```

## Where to get icons

- https://selfh.st/icons/

## Helm

TODO

## ArgoCD

TODO

# Inspiration

- Homer (https://github.com/bastienwirtz/homer)

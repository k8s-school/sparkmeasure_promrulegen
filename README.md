# sparkmeasure_promrulegen

Generate prometheus rule from sparkmeasure metrics

```bash
go build -o rulegen ./cmd/rulegen
cat metrics.txt | ./rulegen > rules.yaml
```
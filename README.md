# ood-stepfunctions-adapter

An [Open OnDemand](https://openondemand.org/) compute adapter that translates OOD job submissions into [AWS Step Functions](https://aws.amazon.com/step-functions/) executions.

## Commands

| Command | Description |
|---------|-------------|
| `submit` | Read a JSON job spec from stdin and start a Step Functions execution; prints the execution ARN |
| `status <execution-arn>` | Return OOD-normalized job status as JSON |
| `delete <execution-arn>` | Stop a running execution (maps to OOD job cancel) |
| `info <execution-arn>` | Print full execution detail as JSON |

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--region` | `us-east-1` | AWS region |
| `--state-machine-arn` | *(required)* | ARN of the target state machine (can also be set per job via `state_machine_arn` in the job spec) |

## Job Spec (stdin for `submit`)

```json
{
  "state_machine_arn": "arn:aws:states:us-east-1:123456789012:stateMachine:MyWorkflow",
  "job_name": "my-run-001",
  "input": {
    "sample": "sample-42",
    "reference": "hg38"
  }
}
```

`state_machine_arn` in the job spec is overridden by `--state-machine-arn` if both are provided.

## OOD Cluster YAML

```yaml
v2:
  metadata:
    title: "AWS Step Functions"
  login:
    host: localhost
  job:
    adapter: linux_host
    submit_host: localhost
    submit:
      launcher: /usr/local/bin/ood-stepfunctions-adapter submit --region us-east-1 --state-machine-arn arn:aws:states:us-east-1:123456789012:stateMachine:MyWorkflow
    status:
      launcher: /usr/local/bin/ood-stepfunctions-adapter status
    delete:
      launcher: /usr/local/bin/ood-stepfunctions-adapter delete
```

## AWS Credentials

The adapter uses the standard AWS credential chain (environment variables, `~/.aws/credentials`, IAM instance profile, etc.). The IAM principal needs `states:StartExecution`, `states:DescribeExecution`, and `states:StopExecution` on the target state machine.

## License

MIT © 2026 Scott Friedman

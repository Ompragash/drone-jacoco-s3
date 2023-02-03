# jacoco-s3-upload-publish

Custom Drone plugin to upload jacoco code coverage reports to AWS S3 bucket and publish the bucket static site url to `Harness CI Execution > Artifacts` tab.

## Build

Build the binary with the following commands:

```bash
go build
```

## Docker

Build the Docker image with the following commands:

```
./hacking/build.sh
docker buildx build -t DOCKER_ORG/drone-jacoco-s3-upload-publish --platform linux/amd64 .
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-s3' not found or does not exist..
```

## Usage

```bash
docker run --rm \
  -e PLUGIN_PIPELINE_SID=44 \
  -e PLUGIN_AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
  -e PLUGIN_AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
  -e PLUGIN_AWS_DEFAULT_REGION=ap-southeast-2 \
  -e PLUGIN_AWS_BUCKET=bucket-name \
  -e PLUGIN_REPORT_SOURCE=maven-code-coverage/target/site/jacoco \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  ompragash/drone-jacoco-s3-upload-publish:latest
```
> _NOTE: Harness Built-in variable `<+pipeline.sequenceId>` can be used as a input value for PLUGIN_PIPELINE_SID in CI pipeline Plugin Step_

In Harness CI,
```yaml
              - step:
                  type: Plugin
                  name: Publish Jacoco Metadata
                  identifier: custom_plugin
                  spec:
                    connectorRef: account.harnessImage
                    image: ompragash/drone-jacoco-s3-upload-publish:latest
                    settings:
                      pipeline_sid: <+pipeline.sequenceId>
                      aws_access_key_id: <+pipeline.variables.AWS_ACCESS_KEY_ID>
                      aws_secret_access_key: <+pipeline.variables.AWS_SECRET_ACCESS_KEY>
                      aws_default_region: ap-southeast-2
                      aws_bucket: bucket-name
                      report_source: maven-code-coverage/target/site/jacoco
```
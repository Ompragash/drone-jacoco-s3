package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	pluginVersion = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-jacoco-s3-upload-publish"
	app.Usage = "Custom Drone plugin to upload jacoco code coverage reports to AWS S3 bucket and publish the bucket static site url to 'Harness CI Execution > Artifacts' tab."
	app.Action = run
	app.Version = pluginVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "pipeline-sequence-id",
			Usage:  "Harness CIE Pipeline Sequence ID",
			EnvVar: "PLUGIN_PIPELINE_SID",
		},
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key ID",
			EnvVar: "PLUGIN_AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Access Key",
			EnvVar: "PLUGIN_AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-default-region",
			Usage:  "AWS Default Region",
			EnvVar: "PLUGIN_AWS_DEFAULT_REGION",
		},
		cli.StringFlag{
			Name:   "aws-bucket",
			Usage:  "AWS Default Region",
			EnvVar: "PLUGIN_AWS_BUCKET",
		},
		cli.StringFlag{
			Name:   "report-source",
			Usage:  "AWS Default Region",
			EnvVar: "PLUGIN_REPORT_SOURCE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	pipelineSeqID := c.String("pipeline-sequence-id")
	awsAccessKey := c.String("aws-access-key")
	awsSecretKey := c.String("aws-secret-key")
	awsDefaultRegion := c.String("aws-default-region")
	awsBucket := c.String("aws-bucket")
	reportSource := c.String("report-source")
	newFolder := "build-" + pipelineSeqID

	// AWS config commands to set ACCESS_KEY_ID and SECRET_ACCESS_KEY
	exec.Command("aws", "configure", "set", "aws_access_key_id", awsAccessKey).Run()
	exec.Command("aws", "configure", "set", "aws_secret_access_key", awsSecretKey).Run()
	s3cmd := exec.Command("aws", "s3", "cp", reportSource, "s3://"+awsBucket+"/"+newFolder, "--region", awsDefaultRegion, "--recursive")

	out, err := s3cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("Output: %s\n", out)
	// End of S3 upload operation

	urls := "http://" + awsBucket + ".s3-website." + awsDefaultRegion + ".amazonaws.com/" + newFolder + "/index.html"
	artifactFilePath := "artifact.txt"

	files := make([]File, 0)
	files = append(files, File{Name: artifactFilePath, URL: urls})

	return writeArtifactFile(files, artifactFilePath)
}

package integron

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/ganvoa/biopipe-tools/internal"
)

type IntegronFinder struct {
	outputDir string
	logger    internal.Logger
}

func NewIntegronFinder(outputDir string, logger internal.Logger) IntegronFinder {
	ifind := IntegronFinder{}
	ifind.outputDir = outputDir
	ifind.logger = logger
	return ifind
}

func (ifind IntegronFinder) Run(downloadDir string, fastaFile string) (string, error) {

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	imageName := "gempasteur/integron_finder:2.0rc10"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)

	fastaFilePath := fmt.Sprintf(`/%s/%s`, downloadDir, fastaFile)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"--local-max", "--split-results", "--func-annot", "--outdir", ifind.outputDir, fastaFilePath},
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf(`%s/%s`, pwd, downloadDir),
				Target: fmt.Sprintf(`/%s`, downloadDir),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf(`%s/%s`, pwd, ifind.outputDir),
				Target: ifind.outputDir,
			},
		},
	}, nil, nil, "")

	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	ifind.logger.Info("integron finder started")

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	ifind.logger.Info("integron finder finished")
	out, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	buf.ReadFrom(out)

	outString := buf.String()
	ifind.logger.Debug(outString)

	fastaName := strings.TrimSuffix(fastaFile, filepath.Ext(fastaFile))
	IntegronFinderOutputDir := fmt.Sprintf(`%s/%s/Results_Integron_Finder_%s`, pwd, ifind.outputDir, fastaName)

	if _, err := os.Stat(IntegronFinderOutputDir); os.IsNotExist(err) {
		return "", errors.New("cant find integron finder result folder")
	}

	ifind.logger.Infof("integron finder result folder succeed %s", IntegronFinderOutputDir)

	return IntegronFinderOutputDir, nil
}

package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
)

// CmdInit handles the init command
func CmdInit(c *cli.Context) error {
	provider := c.String("provider")
	template := c.String("template")
	output := c.String("output")
	force := c.Bool("force")

	// Determine output file
	if output == "" {
		switch provider {
		case "github":
			output = ".github/workflows/ci.yml"
		case "gitlab":
			output = ".gitlab-ci.yml"
		case "bitbucket":
			output = "bitbucket-pipelines.yml"
		case "azure":
			output = "azure-pipelines.yml"
		default:
			output = ".github/workflows/ci.yml"
		}
	}

	// Check if file exists
	if _, err := os.Stat(output); err == nil && !force {
		return fmt.Errorf("file %s already exists. Use --force to overwrite", output)
	}

	// Create directory if needed
	dir := filepath.Dir(output)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate pipeline content
	content := generatePipelineTemplate(provider, template)

	// Write file
	if err := os.WriteFile(output, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", output, err)
	}

	fmt.Printf("✓ Created %s pipeline: %s\n", provider, output)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Review and customize the pipeline\n")
	fmt.Printf("  2. Test locally: git-ci run -f %s\n", output)
	fmt.Printf("  3. Commit and push to repository\n")

	return nil
}

// generatePipelineTemplate generates a pipeline template
func generatePipelineTemplate(provider, template string) string {
	switch provider {
	case "github":
		return generateGitHubTemplate(template)
	case "gitlab":
		return generateGitLabTemplate(template)
	case "bitbucket":
		return generateBitbucketTemplate(template)
	case "azure":
		return generateAzureTemplate(template)
	default:
		return generateGitHubTemplate(template)
	}
}

// generateGitHubTemplate generates GitHub Actions template
func generateGitHubTemplate(template string) string {
	switch template {
	case "node":
		return githubNodeTemplate
	case "python":
		return githubPythonTemplate
	case "go":
		return githubGoTemplate
	case "docker":
		return githubDockerTemplate
	default:
		return githubBasicTemplate
	}
}

// generateGitLabTemplate generates GitLab CI template
func generateGitLabTemplate(template string) string {
	switch template {
	case "node":
		return gitlabNodeTemplate
	case "python":
		return gitlabPythonTemplate
	case "go":
		return gitlabGoTemplate
	case "docker":
		return gitlabDockerTemplate
	default:
		return gitlabBasicTemplate
	}
}

// generateBitbucketTemplate generates Bitbucket Pipelines template
func generateBitbucketTemplate(template string) string {
	// Implement Bitbucket templates
	return bitbucketBasicTemplate
}

// generateAzureTemplate generates Azure Pipelines template
func generateAzureTemplate(template string) string {
	// Implement Azure templates
	return azureBasicTemplate
}

// Template definitions

const githubBasicTemplate = `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Run tests
      run: echo "Add your test commands here"

  build:
    runs-on: ubuntu-latest
    needs: test

    steps:
    - uses: actions/checkout@v3

    - name: Build
      run: echo "Add your build commands here"
`

const githubNodeTemplate = `name: Node.js CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        node-version: [16.x, 18.x, 20.x]

    steps:
    - uses: actions/checkout@v3

    - name: Use Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v3
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'

    - name: Install dependencies
      run: npm ci

    - name: Run tests
      run: npm test

    - name: Run linter
      run: npm run lint

  build:
    runs-on: ubuntu-latest
    needs: test

    steps:
    - uses: actions/checkout@v3

    - name: Use Node.js
      uses: actions/setup-node@v3
      with:
        node-version: 18.x
        cache: 'npm'

    - name: Install dependencies
      run: npm ci

    - name: Build
      run: npm run build

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: build
        path: dist/
`

const githubPythonTemplate = `name: Python CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        python-version: ["3.8", "3.9", "3.10", "3.11"]

    steps:
    - uses: actions/checkout@v3

    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}

    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -r requirements.txt
        pip install pytest flake8

    - name: Lint with flake8
      run: |
        flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
        flake8 . --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics

    - name: Test with pytest
      run: pytest
`

const githubGoTemplate = `name: Go CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.21

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Run vet
      run: go vet ./...

    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    runs-on: ubuntu-latest
    needs: test

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.21

    - name: Build
      run: go build -v ./...
`

const githubDockerTemplate = `name: Docker CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max
`

const gitlabBasicTemplate = `stages:
  - test
  - build
  - deploy

variables:
  CI: "true"

test:
  stage: test
  script:
    - echo "Running tests..."
    - echo "Add your test commands here"

build:
  stage: build
  script:
    - echo "Building application..."
    - echo "Add your build commands here"
  dependencies:
    - test

deploy:
  stage: deploy
  script:
    - echo "Deploying application..."
    - echo "Add your deployment commands here"
  only:
    - main
  when: manual
`

const gitlabNodeTemplate = `image: node:18

stages:
  - test
  - build
  - deploy

variables:
  CI: "true"

cache:
  paths:
    - node_modules/

before_script:
  - npm ci

test:
  stage: test
  script:
    - npm run test
    - npm run lint
  coverage: '/Lines\s+:\s+(\d+\.\d+)%/'

build:
  stage: build
  script:
    - npm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 week
  dependencies:
    - test
`

const gitlabPythonTemplate = `image: python:3.11

stages:
  - test
  - build
  - deploy

variables:
  PIP_CACHE_DIR: "$CI_PROJECT_DIR/.cache/pip"

cache:
  paths:
    - .cache/pip
    - venv/

before_script:
  - python -m venv venv
  - source venv/bin/activate
  - pip install -r requirements.txt

test:
  stage: test
  script:
    - pip install pytest flake8
    - flake8 .
    - pytest --cov=.
  coverage: '/TOTAL.*\s+(\d+)%/'

build:
  stage: build
  script:
    - python setup.py bdist_wheel
  artifacts:
    paths:
      - dist/
  dependencies:
    - test
`

const gitlabGoTemplate = `image: golang:1.21

stages:
  - test
  - build

variables:
  GO111MODULE: "on"

cache:
  paths:
    - .go/pkg/mod/

before_script:
  - mkdir -p .go
  - export GOPATH=$CI_PROJECT_DIR/.go
  - export PATH=$PATH:$GOPATH/bin

test:
  stage: test
  script:
    - go mod download
    - go test -v -race -coverprofile=coverage.out ./...
    - go vet ./...
    - go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    - golangci-lint run
  coverage: '/coverage: \d+.\d+% of statements/'

build:
  stage: build
  script:
    - go build -v -o app ./cmd/...
  artifacts:
    paths:
      - app
  dependencies:
    - test
`

const gitlabDockerTemplate = `image: docker:latest

services:
  - docker:dind

stages:
  - build
  - push

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"
  IMAGE_TAG: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA

build:
  stage: build
  script:
    - docker build -t $IMAGE_TAG .
    - docker save $IMAGE_TAG > image.tar
  artifacts:
    paths:
      - image.tar
    expire_in: 1 hour

push:
  stage: push
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker load < image.tar
    - docker push $IMAGE_TAG
    - docker tag $IMAGE_TAG $CI_REGISTRY_IMAGE:latest
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main
`

const bitbucketBasicTemplate = `pipelines:
  default:
    - step:
        name: Test
        script:
          - echo "Running tests..."
    - step:
        name: Build
        script:
          - echo "Building application..."
`

const azureBasicTemplate = `trigger:
- main

pool:
  vmImage: ubuntu-latest

stages:
- stage: Test
  jobs:
  - job: Test
    steps:
    - script: echo "Running tests..."
      displayName: 'Run tests'

- stage: Build
  dependsOn: Test
  jobs:
  - job: Build
    steps:
    - script: echo "Building application..."
      displayName: 'Build application'
`

pipeline {
    agent none // We will define the agent for each stage

    stages {
        stage('Test Marketplace') {
            agent {
                docker {
                    // Jenkins will spin up a Go container just for this step!
                    image 'golang:1.23-alpine' 
                    // Mount the current folder so Go can see the code
                    reuseNode true
                }
            }
            steps {
                echo 'Running Tests for Marketplace Service...'
                // 1. Go into the folder
                // 2. Download dependencies
                // 3. Run tests (If no test files exist, this passes harmlessly)
                sh 'cd marketplace-service && go mod download && go test -v ./...'
            }
        }

        stage('Test Watcher') {
            agent {
                docker { image 'golang:1.23-alpine' }
            }
            steps {
                echo 'Running Tests for Watcher Service...'
                sh 'cd watcher-service && go mod download && go test -v ./...'
            }
        }
    }
}
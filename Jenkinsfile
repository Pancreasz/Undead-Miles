pipeline {
    agent none 

    stages {
        stage('Test Marketplace') {
            agent {
                docker {
                    image 'golang:1.24-alpine' 
                    reuseNode true
                }
            }
            steps {
                echo 'Running Tests for Marketplace Service...'
                sh 'cd marketplace-service && go mod download && go test -v ./...'
            }
        }

        stage('Test Watcher') {
            agent {
                docker { image 'golang:1.24-alpine' }
            }
            steps {
                echo 'Running Tests for Watcher Service...'
                sh 'cd watcher-service && go mod download && go test -v ./...'
            }
        }
    }
}
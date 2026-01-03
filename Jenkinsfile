pipeline {
    agent none 

    stages {
        // --- STAGE 1: TEST ---
        stage('Test Services') {
            agent {
                docker { 
                    image 'golang:1.24-alpine' // Using your new version
                    reuseNode true
                }
            }
            steps {
                echo 'Testing Marketplace...'
                sh 'cd marketplace-service && go mod download && go test -v ./...'
                
                echo 'Testing Watcher...'
                sh 'cd watcher-service && go mod download && go test -v ./...'
            }
        }

        // --- STAGE 2: BUILD (New!) ---
        stage('Build Docker Images') {
            agent any
            steps {
                echo 'Building Docker Images...'
                script {
                    // This command uses the host's Docker to build the image
                    // We tag it with 'jenkins-built' so you can easily find it later
                    sh 'docker build -t undeadmiles-marketplace:jenkins-built ./marketplace-service'
                    sh 'docker build -t undeadmiles-watcher:jenkins-built ./watcher-service'
                }
            }
        }
    }
}
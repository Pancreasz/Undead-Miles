pipeline {
    agent {
        kubernetes {
            // We define a "Pod" that has all the tools we need inside it
            yaml '''
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              # Container 1: Golang for testing
              - name: golang
                image: golang:1.24-alpine
                command:
                - sleep
                - infinity
                
              # Container 2: Kaniko for building images (replaces 'docker build')
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command:
                - sleep
                - infinity
                env:
                  - name: PATH
                    value: /usr/local/bin:/kaniko
            '''
        }
    }

    stages {
        // --- STAGE 1: TEST (Runs in Golang container) ---
        stage('Test Services') {
            steps {
                container('golang') {
                    echo 'Testing Marketplace...'
                    sh 'cd marketplace-service && go mod download && go test -v ./...'
                    
                    echo 'Testing Watcher...'
                    sh 'cd watcher-service && go mod download && go test -v ./...'
                }
            }
        }

        // --- STAGE 2: BUILD (Runs in Kaniko container) ---
        stage('Build Docker Images') {
            steps {
                container('kaniko') {
                    echo 'Building Marketplace Image...'
                    // Kaniko builds the image without needing a Docker daemon
                    // --no-push ensures we just test the build (since we haven't set up registry credentials yet)
                    sh '''
                    /kaniko/executor --context ./marketplace-service \
                    --dockerfile ./marketplace-service/Dockerfile \
                    --destination my-repo/undeadmiles-marketplace:jenkins-built \
                    --no-push 
                    '''

                    echo 'Building Watcher Image...'
                    sh '''
                    /kaniko/executor --context ./watcher-service \
                    --dockerfile ./watcher-service/Dockerfile \
                    --destination my-repo/undeadmiles-watcher:jenkins-built \
                    --no-push
                    '''
                }
            }
        }
    }
}
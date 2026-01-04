pipeline {
    agent {
        kubernetes {
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
                
              # Container 2: Kaniko for building images
              # We MUST use the :debug tag to get a shell
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                # We use /busybox/cat to keep the container running since standard 'sleep' is missing
                command:
                - /busybox/cat
                tty: true
                volumeMounts:
                  - name: kaniko-secret
                    mountPath: /kaniko/.docker
              volumes:
                - name: kaniko-secret
                  secret:
                    secretName: docker-registry-creds
                    items:
                      - key: .dockerconfigjson
                        path: config.json
                    optional: true # Make optional for now so it runs even without creds
            '''
        }
    }

    stages {
        // --- STAGE 1: TEST ---
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

        // --- STAGE 2: BUILD ---
        stage('Build Docker Images') {
            steps {
                container('kaniko') {
                    // We must tell Jenkins to use the busybox shell
                    script {
                        echo 'Building Marketplace Image...'
                        sh '''#!/busybox/sh
                        /kaniko/executor --context ./marketplace-service \
                        --dockerfile ./marketplace-service/Dockerfile \
                        --destination my-repo/undeadmiles-marketplace:jenkins-built \
                        --no-push
                        '''

                        echo 'Building Watcher Image...'
                        sh '''#!/busybox/sh
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
}
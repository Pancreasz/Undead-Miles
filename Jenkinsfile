pipeline {
    agent {
        kubernetes {
            yaml '''
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              # 1. Golang Container
              - name: golang
                image: golang:1.24-alpine
                command: ["sleep", "infinity"]
                
              # 2. Kaniko Container
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command: ["/busybox/cat"]
                tty: true
                # === FIX 1: Explicitly set the working directory ===
                workingDir: /home/jenkins/agent
                resources:
                  limits:
                    memory: "2Gi"
                    cpu: "1"
                  requests:
                    memory: "512Mi"
                    cpu: "500m"
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
                    optional: true
            '''
        }
    }

    stages {
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

        stage('Build Docker Images') {
            steps {
                container('kaniko') {
                    // === FIX 2: Added explicit 'dir' block to be safe ===
                    // This ensures Jenkins runs these commands in the root of your project
                    dir('/home/jenkins/agent/workspace/UndeadMiles-Pipeline') {
                        
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
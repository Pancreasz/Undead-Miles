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
                
              # 2. Kaniko Container (UPDATED WITH RESOURCES)
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command: ["/busybox/cat"]
                tty: true
                # === NEW: Give Kaniko enough RAM to compile Go ===
                resources:
                  requests:
                    memory: "512Mi"
                    cpu: "500m"
                  limits:
                    memory: "2Gi"  # 2GB should be plenty for Go builds
                    cpu: "1"
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
                    // Tip: Split these into separate steps so logs are clearer
                    // and Jenkins can manage resources better between them.
                    
                    echo 'Building Marketplace Image...'
                    sh '''#!/busybox/sh
                    /kaniko/executor --context ./marketplace-service \
                    --dockerfile ./marketplace-service/Dockerfile \
                    --destination my-repo/undeadmiles-marketplace:jenkins-built \
                    --no-push
                    '''

                    echo 'Cleaning up before next build...'
                    // Optional: Brief pause or cleanup if needed, but RAM increase usually fixes it.

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
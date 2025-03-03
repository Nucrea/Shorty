pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'docker build . -t shorty'
            }
        }
        stage('Deploy') {
            steps {
                sh 'docker stop shorty || true'
                sh 'docker rm shorty || true'
                sh 'docker run --name=shorty ${SHORTY_APP_PORT}:8081 -e SHORTY_POSTGRES_URL=${SHORTY_POSTGRES_URL} -e SHORTY_BASE_URL=${SHORTY_BASE_URL} -d shorty'
            }
        }
    }
}
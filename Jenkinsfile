
pipeline {

    agent none

    environment {
        GITHUB_TOKEN = credentials('GITHUB_TOKEN')
    }

    stages {

        stage('Build & Test') {
            agent { dockerfile true}
            steps {
                echo "Building via docker"
            }
        }
    }

    post {
        failure {
            emailext attachLog: true, body: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS:\n Check console output at $BUILD_URL to view the results.\n\n', recipientProviders: [culprits()], subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!'
        }
    }
}

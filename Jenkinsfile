
pipeline {

    environment {
        GITHUB_TOKEN = credentials('GITHUB_TOKEN')
    }

    agent { dockerfile }

    stages {

        stage('Build & Test') {
            steps {
                agent { dockerfile true }
            }
        }
    }

    post {
        failure {
            emailext attachLog: true, body: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS:\n Check console output at $BUILD_URL to view the results.\n\n', recipientProviders: [culprits()], subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!'
        }
    }
}

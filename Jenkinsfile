
pipeline {

    agent any

    environment {
        GITHUB_TOKEN = credentials('GITHUB_TOKEN')
        DOCKER_REPOSITORY = "docker.fg"
    }

    stages {

        stage('Build docker image') {
            steps {
                echo "Building image..."
                sh "docker build -t $DOCKER_REPOSITORY/fg-vflow:$BUILD_NUMBER"
                echo "Build image complete"
            }
        }

        stage ('Push docker image') {
            steps {
                echo "Pushing docker image..."
                sh "docker push $DOCKER_REPOSITORY/fg-vflow:$BUILD_NUMBER"
                echo "Push image complete"
            }
        }
    }

    post {
        success {
            echo "Success"
        }
        failure {
            emailext attachLog: true, body: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS:\n Check console output at $BUILD_URL to view the results.\n\n', recipientProviders: [culprits()], subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!'
        }
    }
}

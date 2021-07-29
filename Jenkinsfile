
pipeline {
    agent any
    environment {
        GITHUB_TOKEN = credentials('GITHUB_TOKEN')
        DOCKER_REPOSITORY = "docker.fg"
        GIT_VERSION=sh(script: 'git describe --tags --always', returnStdout: true).toString().trim()
    }

    stages {

        stage ("Preparing container for golang")  {
            agent {
                docker {
                    image "golang"
                }
            }

            stages {
                stage ("Unit testing") {
                    steps {
                        echo "Unit testing..."
                        sh "go build ./..."
                        sh "go install ./..."
                        sh "go test -v ./... -timeout 1m"
                    }
                }
            }
        }

        stage ("Code quality") {
            steps {
                echo("Checking code quality....")
                script {
                    def scannerHome = tool 'Sonar Scanner 3.0.0.702';
                    withSonarQubeEnv {
                        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectVersion=$GIT_VERSION"
                    }
                }
            }
        }


        stage("Docker build & publish") {
            steps {
                script {
                    dockerImage = docker.build "$DOCKER_REPOSITORY/fg_vflow"

                    bn = env.BUILD_NUMBER
                    currentBuild.displayName = "#${bn}:$GIT_VERSION"

                    dockerImage.push($GIT_VERSION)
                    if (env.BRANCH_NAME == "devel") {
                        dockerImage.push("devel")
                    }
                }
            }
        }

        stage ("Devel deploy") {
            when { branch "devel" }
            steps {
                salt(authtype: 'pam',
                    clientInterface: local(arguments: 'node.rtbh', blockbuild: true, function: 'state.apply', jobPollTime: 6, target: 'node-1.bohdalec.test.fg', targettype: 'glob'),
                    credentialsId: '3f36bac7-b50e-42f2-b977-19e352fbd3c7',
                    saveFile: true,
                    servername: 'https://salt.test.fg:8000/')
                script {
                    env.WORKSPACE = pwd()
                    def output = readFile "${env.WORKSPACE}/saltOutput.json"
                    echo output
                    echo "Done..."
                }
            }
        }
    }

    post {
        success {
            echo "Pipeline success"
        }
        failure {
            echo "Pipeline failed..."
            emailext attachLog: true, body: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS:\n Check console output at $BUILD_URL to view the results.\n\n', recipientProviders: [culprits()], subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!'
        }
    }
}

node {
    // Install the desired Go version
    def ROOT = tool name: 'Go 1.12.1', type: 'go'

    // Variables
    def DOCKER_IMAGE

    def APP_NAME = "golang-rh-todo"
    def USERNAME = "karolispx"

    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        withEnv(["GOROOT=${ROOT}", "PATH+GO=${ROOT}/bin"]) {

            // Clone The Repository
            stage('Clone The Repository') {
                sh "git clone https://github.com/${USERNAME}/${APP_NAME}.git"
                // This clones everything into folder specified in APP_NAME variable
                // Since this is a public repository, we don't have to worry about logins or keys
            }

            // Install Dependencies
            stage('Install Dependencies') {
                // Make sure to switch to appropriare folder specified in APP_NAME variable
                dir("${APP_NAME}") {
                    sh "go get github.com/gorilla/handlers"
                    sh "go get github.com/gorilla/mux"
                    sh "go get github.com/karolispx/golang-rh-todo/api"
                    sh "go get github.com/karolispx/golang-rh-todo/helpers"
                    sh "go get github.com/karolispx/golang-rh-todo/models"
                }
            }

            // Run Tests
            stage('Run Tests') {
                // Make sure to switch to appropriare folder specified in APP_NAME variable
                dir("${APP_NAME}") {
                    echo "This is where I would run tests, if I had some!"
                }
            }

            // Build Binary from Go
            stage('Build Binary from Go') {
                // Make sure to switch to appropriare folder specified in APP_NAME variable
                dir("${APP_NAME}") {
                    sh "env GOOS=linux GOARCH=amd64 go build -o bin/${APP_NAME}"
                }
            }

            // Build Docker Image - it will use Dockerfile that's currently in the root directory of the project
            stage('Build Docker Image') {
                dir("${APP_NAME}") {
                    DOCKER_IMAGE = docker.build("${USERNAME}/${APP_NAME}")
                }
            }

            // Push Docker Image - this will push it to docker hub.
            stage('Push Docker Image') {
                dir("${APP_NAME}") {
                    docker.withRegistry('https://registry.hub.docker.com', 'DOCKER-HUB-CREDENTIALS') {
                        DOCKER_IMAGE.push("latest")
                    }
                }
            }

            // Finished
            stage('Finished') {
                echo "Finished successfully!"
            }
        }
    }
}
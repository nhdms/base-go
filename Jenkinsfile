pipeline {
    agent any

    environment {
        // Core configuration
        DOCKER_TAG = "${env.GIT_COMMIT[0..7]}-${env.BUILD_NUMBER}"
        // Repository settings
        DEPLOYMENT_REPO = 'git@github.com/deployments.git'
        DEPLOYMENT_REPO_ID = '64636931'

        // ArgoCD configuration
        ARGOCD_CONFIG = getArgoConfig(BRANCH_NAME)
        ARGOCD_SERVER = "${ARGOCD_CONFIG.server}"
        ARGOCD_APP = "${ARGOCD_CONFIG.app}"
        ARGOCD_CREDENTIALS = "${ARGOCD_CONFIG.credentials}"

        AUTO_MERGE = 'true'
        SERVICE_PATH = extractServicePath(env.JOB_NAME)
    }

    stages {
        stage('Prepare Environment') {
            steps {
                script {
                    def branchName = env.BRANCH_NAME
                    def cicdPath = "${env.SERVICE_PATH}/config/cicd.json"
                    // Only proceed with auto-merge if branch is 'dev'
                    env.AUTO_MERGE = (branchName == 'dev') ? 'true' : 'false'
                    if (!fileExists(cicdPath)) {
                        error "Deployment file not found: ${cicdPath}"
                    }

                    def cicdConfig = readJSON file: cicdPath
                    
                    // Set SERVICE_NAME using parsed JSON
                    env.DOCKER_IMAGE = "agbiz/go-${cicdConfig.app_type}-${cicdConfig.service_name}"
                    env.SERVICE_TYPE = "${cicdConfig.app_type}"
                    env.SERVICE_NAME = "${cicdConfig.service_name}"
                    env.MAPPED_SERVICE_NAME = "${cicdConfig.service_name}-${cicdConfig.app_type}-app"
                    env.DEPLOY_BRANCH_NAME = "update-image-${env.SERVICE_NAME}-${env.BUILD_NUMBER}"

                    echo "Branch name: ${branchName}"
                    echo "CICD config path: ${cicdPath}"
                    echo "Service name: ${env.SERVICE_NAME}"
                    echo "Image name: ${env.DOCKER_IMAGE}"
                }
            }
        }

        stage('Build and Push Docker Image') {
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'docker-hub-credentials',
                               usernameVariable: 'DOCKER_USERNAME',
                                passwordVariable: 'DOCKER_PASSWORD'),
                               sshUserPrivateKey(credentialsId: 'gitlab-ssh-key', 
                                    keyFileVariable: 'SSH_KEY_PATH'),
                               ]) {
                    // Build and push image
                    sh """
                        cp "\${SSH_KEY_PATH}" ./id_rsa
                        chmod 600 ./id_rsa
                        DOCKER_BUILDKIT=1 docker build --platform linux/amd64 --build-arg BIN=${SERVICE_PATH} -f scripts/Dockerfile --ssh default=./id_rsa . -t ${DOCKER_IMAGE}:${DOCKER_TAG}
                        echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
                        docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                    """
                }
            }
        }

        stage('Update Deployment') {
            steps {
                sshagent(['gitlab-ssh-key']) {
                    script {
                        // Clone and prepare repository
                        sh """
                            mkdir -p ~/.ssh
                            ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
                            rm -rf ${MAPPED_SERVICE_NAME}
                            git clone ${DEPLOYMENT_REPO} ${MAPPED_SERVICE_NAME}
                        """

                        dir(MAPPED_SERVICE_NAME) {
                            // Update deployment file
                            def deploymentFile = "overlays/${BRANCH_NAME}/${SERVICE_TYPE}/${SERVICE_NAME}/kustomization.yaml"
                            if (!fileExists(deploymentFile)) {
                                error "Deployment file not found: ${deploymentFile}"
                            }

                            // Update image and create commit
                            sh """
                                git checkout -b ${DEPLOY_BRANCH_NAME}
                                # Update newTag in kustomization.yaml using sed
                                sed -i "s/newTag: .*/newTag: '${env.DOCKER_TAG}'/" ${deploymentFile}
                                
                                # Verify the change
                                cat ${deploymentFile}

                                git add ${deploymentFile}
                                git commit -m 'Update ${MAPPED_SERVICE_NAME} image to ${DOCKER_TAG}'
                                git push -uf origin ${DEPLOY_BRANCH_NAME}
                            """
                        }
                    }
                }
            }
        }

        stage('Create and Process Merge Request') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'gitlab-credentials',
                            usernameVariable: 'USERNAME', passwordVariable: 'GITLAB_TOKEN')]) {
                    script {
                        // Create merge request
                        def mrId = createMergeRequest(GITLAB_TOKEN)
                        echo "Created merge request #${mrId}"

                        if (env.AUTO_MERGE == 'true') {
                            // Wait for merge request to be ready with timeout
                            def startTime = System.currentTimeMillis()
                            def timeoutMillis = 30000 // 30 seconds
                            def isReady = false

                            while (System.currentTimeMillis() - startTime < timeoutMillis) {
                                def mrStatus = sh(
                                    script: """
                                        curl --silent --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
                                            https://gitlab.com/api/v4/projects/${env.DEPLOYMENT_REPO_ID}/merge_requests/${mrId} \
                                            | jq -r '.merge_status'
                                    """,
                                    returnStdout: true
                                ).trim()

                                if (mrStatus == 'can_be_merged') {
                                    isReady = true
                                    mergeMergeRequest(GITLAB_TOKEN, mrId)
                                    echo "Merge request #${mrId} successfully merged"
                                    break
                                }

                                sleep(5) // Wait 5 seconds before next check
                            }

                            if (!isReady) {
                                error "Timeout: Merge request #${mrId} was not ready to merge within 30 seconds"
                            }
                        }
                    }
                }
            }
        }

        stage('Sync ArgoCD') {
            steps {
                withCredentials([usernamePassword(credentialsId: "${ARGOCD_CREDENTIALS}",
                               usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                    sh """
                        argocd login ${ARGOCD_SERVER} \
                            --username ${USERNAME} \
                            --password ${PASSWORD} \
                            --grpc-web
                        argocd app sync ${ARGOCD_APP} \
                            --resource apps:Deployment:${MAPPED_SERVICE_NAME} \
                            --grpc-web
                    """
                }
            }
        }
    }
}

def createMergeRequest(token) {
    def response = sh(
        script: """
            curl --silent --header "PRIVATE-TOKEN: ${token}" \
                --data "remove_source_branch=true&source_branch=${env.DEPLOY_BRANCH_NAME}&target_branch=main&title=Update ${env.MAPPED_SERVICE_NAME} image to ${DOCKER_TAG}" \
                https://gitlab.com/api/v4/projects/${env.DEPLOYMENT_REPO_ID}/merge_requests \
                | jq -r '.iid'
        """,
        returnStdout: true
    ).trim()
    return response
}

def mergeMergeRequest(token, mrId) {
    sh """
        curl --silent --header "PRIVATE-TOKEN: ${token}" \
            --request PUT \
            https://gitlab.com/api/v4/projects/${env.DEPLOYMENT_REPO_ID}/merge_requests/${mrId}/merge
    """
}

def extractServicePath(jobName) {
    try {
        // Get the repository name from job path (second segment)
        def repoName = jobName.tokenize('/')[1]

        // Split by first underscore: "consumer_folder1_folder2_service-name"
        def firstSplit = repoName.split('_', 2)
        if (firstSplit.size() != 2) {
            error "Invalid job name format. Expected: type_path1_path2_service-name, got: ${repoName}"
        }

        def serviceType = firstSplit[0]     // "consumer"
        def remainingPath = firstSplit[1]   // "folder1_folder2_service-name"

        // Replace all underscores with slashes
        def servicePath = remainingPath.replace('_', '/')

        // Construct the full path with correct pluralization
        return "cmd/${serviceType}s/${servicePath}"
    } catch (Exception e) {
        error "Failed to extract service path from job name: ${jobName}\nError: ${e.message}"
    }
}

def getArgoConfig(branchName) {
    def config = [:]
    config.app = "go-services-${branchName}"

    if (branchName == 'production') {
        config.server = 'server-prod.domain.com'
        config.credentials = 'argo_cd_admin_prod'
        return config
    }
    config.server = 'server-dev.domain.com'
    config.credentials = 'argo_cd_admin'
    return config
}
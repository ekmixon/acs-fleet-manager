die() {
    echo >&2 "$*"
    exit 1
}

log() {
    echo "$*"
}

verify_environment() {
    return
}

get_current_cluster_name() {
    local kubeconfig_file="$1"
    local kubeconf=$(cat "$kubeconfig_file")
    local context_name=$(echo "$kubeconf" | yq e .current-context -)
    local context="$(echo "$kubeconf" | yq e ".contexts[] | select(.name == \"${context_name}\")" -j - | jq -c)"
    local cluster_name=$(echo "$context" | jq -r .context.cluster -)
    echo "$cluster_name"
}

init() {
    export DEBUG="${DEBUG:-false}"
    set -eu -o pipefail
    if [[ "$DEBUG" == "trace" ]]; then
        set -x
    fi

    source "${GITROOT}/dev/env/defaults/env"
    if [[ -n "$OPENSHIFT_CI" ]]; then
        source "${GITROOT}/dev/env/defaults/openshift-ci.env"
    fi

    export CLUSTER_TYPE="${CLUSTER_TYPE:-$CLUSTER_TYPE_DEFAULT}"
    source "${GITROOT}/dev/env/defaults/cluster-type-${CLUSTER_TYPE}/env"

    export ACSMS_NAMESPACE="${ACSMS_NAMESPACE:-$ACSMS_NAMESPACE_DEFAULT}"
    export KUBECTL=${KUBECTL:-$KUBECTL_DEFAULT}
    export DOCKER=${DOCKER:-$DOCKER_DEFAULT}
    export IMAGE_REGISTRY="${IMAGE_REGISTRY:-$IMAGE_REGISTRY_DEFAULT}"
    export STACKROX_OPERATOR_VERSION="${STACKROX_OPERATOR_VERSION:-$STACKROX_OPERATOR_VERSION_DEFAULT}"
    export CENTRAL_VERSION="${CENTRAL_VERSION:-$CENTRAL_VERSION_DEFAULT}"
    export SCANNER_VERSION="${SCANNER_VERSION:-$SCANNER_VERSION_DEFAULT}"
    export STACKROX_OPERATOR_NAMESPACE="${STACKROX_OPERATOR_NAMESPACE:-$STACKROX_OPERATOR_NAMESPACE_DEFAULT}"
    export STACKROX_OPERATOR_IMAGE="${IMAGE_REGISTRY}/stackrox-operator:${STACKROX_OPERATOR_VERSION}"
    export STACKROX_OPERATOR_INDEX_IMAGE="${IMAGE_REGISTRY}/stackrox-operator-index:v${STACKROX_OPERATOR_VERSION}"
    export KUBECONFIG=${KUBECONFIG:-$KUBECONFIG_DEFAULT}
    export CLUSTER_NAME_DEFAULT=$(get_current_cluster_name "$KUBECONFIG")
    export CLUSTER_NAME=${CLUSTER_NAME:-$CLUSTER_NAME_DEFAULT}
    export OPENSHIFT_MARKETPLACE="${OPENSHIFT_MARKETPLACE:-$OPENSHIFT_MARKETPLACE_DEFAULT}"
    export INSTALL_OPERATOR="${INSTALL_OPERATOR:-$INSTALL_OPERATOR_DEFAULT}"
    export POSTGRES_DB=${POSTGRES_DB:-$POSTGRES_DB_DEFAULT}
    export POSTGRES_USER=${POSTGRES_USER:-$POSTGRES_USER_DEFAULT}
    export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-$POSTGRES_PASSWORD_DEFAULT}
    export DATABASE_HOST="db"
    export DATABASE_PORT="5432"
    export DATABASE_NAME="$POSTGRES_DB"
    export DATABASE_USER="$POSTGRES_USER"
    export DATABASE_PASSWORD="$POSTGRES_PASSWORD"
    export DATABASE_TLS_CERT=""
    export OCM_SERVICE_CLIENT_ID=${OCM_SERVICE_CLIENT_ID:-$OCM_SERVICE_CLIENT_ID_DEFAULT}
    export OCM_SERVICE_CLIENT_SECRET=${OCM_SERVICE_CLIENT_SECRET:-$OCM_SERVICE_CLIENT_SECRET_DEFAULT}
    export OCM_SERVICE_TOKEN=${OCM_SERVICE_TOKEN:-$OCM_SERVICE_TOKEN_DEFAULT}
    export SENTRY_KEY=${SENTRY_KEY:-$SENTRY_KEY_DEFAULT}
    export AWS_ACCESS_KEY=${AWS_ACCESS_KEY:-$AWS_ACCESS_KEY_DEFAULT}
    export AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID:-$AWS_ACCOUNT_ID_DEFAULT}
    export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-$AWS_SECRET_ACCESS_KEY_DEFAULT}
    export SSO_CLIENT_ID=${SSO_CLIENT_ID:-$SSO_CLIENT_ID_DEFAULT}
    export SSO_CLIENT_SECRET=${SSO_CLIENT_SECRET:-$SSO_CLIENT_SECRET_DEFAULT}
    export OSD_IDP_SSO_CLIENT_ID=${OSD_IDP_SSO_CLIENT_ID:-$OSD_IDP_SSO_CLIENT_ID_DEFAULT}
    export OSD_IDP_SSO_CLIENT_SECRET=${OSD_IDP_SSO_CLIENT_SECRET:-$OSD_IDP_SSO_CLIENT_SECRET_DEFAULT}
    export ROUTE53_ACCESS_KEY=${ROUTE53_ACCESS_KEY:-$ROUTE53_ACCESS_KEY_DEFAULT}
    export ROUTE53_SECRET_ACCESS_KEY=${ROUTE53_SECRET_ACCESS_KEY:-$ROUTE53_SECRET_ACCESS_KEY_DEFAULT}
    export OBSERVABILITY_CONFIG_ACCESS_TOKEN=${OBSERVABILITY_CONFIG_ACCESS_TOKEN:-$OBSERVABILITY_CONFIG_ACCESS_TOKEN_DEFAULT}
    export IMAGE_PULL_DOCKER_CONFIG=${IMAGE_PULL_DOCKER_CONFIG:-$IMAGE_PULL_DOCKER_CONFIG_DEFAULT}
    export KUBECONF_CLUSTER_SERVER_OVERRIDE=${KUBECONF_CLUSTER_SERVER_OVERRIDE:-$KUBECONF_CLUSTER_SERVER_OVERRIDE_DEFAULT}
    export INHERIT_IMAGEPULLSECRETS=${INHERIT_IMAGEPULLSECRETS:-$INHERIT_IMAGEPULLSECRETS_DEFAULT}
    export SPAWN_LOGGER=${SPAWN_LOGGER:-$SPAWN_LOGGER_DEFAULT}
    export OPERATOR_SOURCE=${OPERATOR_SOURCE:-$OPERATOR_SOURCE_DEFAULT}
    export INSTALL_OLM=${INSTALL_OLM:-$INSTALL_OLM_DEFAULT}

    export FLEET_MANAGER_IMAGE="${FLEET_MANAGER_IMAGE:-$FLEET_MANAGER_IMAGE_DEFAULT}"
    # When transferring images without repository hostname to Minikube it gets prefixed with "docker.io" automatically.
    if [[ "$FLEET_MANAGER_IMAGE" =~ ^fleet-manager-.*/fleet-manager:.* ]]; then
        export FULL_FLEET_MANAGER_IMAGE="docker.io/${FLEET_MANAGER_IMAGE}"
    else
        export FULL_FLEET_MANAGER_IMAGE="${FLEET_MANAGER_IMAGE}"
    fi

    verify_environment

    disable_debugging
    enable_debugging_if_necessary
}

disable_debugging() {
    if [[ "$DEBUG" != "trace" ]]; then
        set +x
    fi
}

enable_debugging_if_necessary() {
    if [[ "$DEBUG" != "false" ]]; then
        set -x
    fi
}
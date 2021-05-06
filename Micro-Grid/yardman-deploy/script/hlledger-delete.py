# -*- coding: utf-8 -*-
import sys
import os
import threading

def delete_command(file_path):
    command = "kubectl delete -f " + file_path + " --grace-period=0 --force"
    # command = "kubectl delete -f " + file_path
    os.system(command)

# delete peer deployment
def delete_deployment(deployment):
    files = os.listdir(deployment)
    for file_name in files:
        if not file_name.startswith("namespace"):
            file_path = deployment + "/" + file_name
            delete_command(file_path)
        else:
            namespace_file = deployment + "/" + file_name
    # delete namespace
    delete_command(namespace_file)
    

if __name__ == "__main__":
    yardman_name = sys.argv[1]
    cluster_name = sys.argv[2]
    if os.path.exists('deploy'):
        deployment = "../deploy/data/" + yardman_name + "/" + cluster_name + "/deployment/"
    else:
        deployment = "../data/" + yardman_name + "/" + cluster_name + "/deployment/"
    if not os.path.exists(deployment):
        print("The cluster is deleted.")
        exit(0)

    orgs = os.listdir(deployment)
    for org in orgs:
        org_deployment = deployment + org
        delete_org = threading.Thread(target=delete_deployment, args=(org_deployment, ))
        delete_org.start()

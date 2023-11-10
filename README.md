## Hackweek 23 : RKE2 and Rancher deployment with Pulumi

### Project Description
Terraform and Ansible are used within the SUSE Consulting Team for automating RKE2 and Rancher deployments, but with the change in Open-Source License of Terraform and the RedHat “problems”, I think there is a need for an alternative solution like Pulumi. I have no experience with Pulumi and there isn’t much documentation around it (regarding Rancher and Terraform). There is a Package for Rancher but nothing for RKE2.

### Goal
Build an example solution for installing RKE2 and Rancher with Pulumi (including things like Longhorn and other Apps) and present it to the team.

### Result
Developing workflows in Pulumi is way more complex that for example Ansible (or even Terraform). The support for system configuration tasks is not existent. There is only a package for SSH command and filecopy. Because the Hackweek has ended, the following is now supported in the example code.

- Creation of KVM VM's with libvirt.
- Installation of RKE2 (non-airgapped and airgapped)

### Resources
- [Hackweek 23 Project](https://hackweek.opensuse.org/23/projects/rke2-and-rancher-deployment-with-pulumi)
- [Pulumi](https://www.pulumi.com)
- [RKE2 Docs](https://docs.rke2.io/)
- [Rancher Docs](https://ranchermanager.docs.rancher.com/)

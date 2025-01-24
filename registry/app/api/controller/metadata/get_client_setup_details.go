//  Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
	"context"
	"net/http"
	"strings"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/app/paths"
	"github.com/harness/gitness/registry/app/api/openapi/contracts/artifact"
	"github.com/harness/gitness/registry/app/common"
	"github.com/harness/gitness/types/enum"
)

func (c *APIController) GetClientSetupDetails(
	ctx context.Context,
	r artifact.GetClientSetupDetailsRequestObject,
) (artifact.GetClientSetupDetailsResponseObject, error) {
	regRefParam := r.RegistryRef
	imageParam := r.Params.Artifact
	tagParam := r.Params.Version

	regInfo, _ := c.GetRegistryRequestBaseInfo(ctx, "", string(regRefParam))

	space, err := c.SpaceStore.FindByRef(ctx, regInfo.ParentRef)
	if err != nil {
		return artifact.GetClientSetupDetails400JSONResponse{
			BadRequestJSONResponse: artifact.BadRequestJSONResponse(
				*GetErrorResponse(http.StatusBadRequest, err.Error()),
			),
		}, nil
	}

	session, _ := request.AuthSessionFrom(ctx)
	permissionChecks := GetPermissionChecks(space, regInfo.RegistryIdentifier, enum.PermissionRegistryView)
	if err = apiauth.CheckRegistry(
		ctx,
		c.Authorizer,
		session,
		permissionChecks...,
	); err != nil {
		return artifact.GetClientSetupDetails403JSONResponse{
			UnauthorizedJSONResponse: artifact.UnauthorizedJSONResponse(
				*GetErrorResponse(http.StatusForbidden, err.Error()),
			),
		}, nil
	}

	reg, err := c.RegistryRepository.GetByParentIDAndName(ctx, regInfo.parentID, regInfo.RegistryIdentifier)
	if err != nil {
		return artifact.GetClientSetupDetails404JSONResponse{
			NotFoundJSONResponse: artifact.NotFoundJSONResponse(
				*GetErrorResponse(http.StatusNotFound, "registry doesn't exist with this ref"),
			),
		}, err
	}

	if imageParam != nil {
		_, err := c.ImageStore.GetByName(ctx, reg.ID, string(*imageParam))
		if err != nil {
			return artifact.GetClientSetupDetails404JSONResponse{
				NotFoundJSONResponse: artifact.NotFoundJSONResponse(
					*GetErrorResponse(http.StatusNotFound, "image doesn't exist"),
				),
			}, err
		}
		if tagParam != nil {
			_, err := c.TagStore.FindTag(ctx, reg.ID, string(*imageParam), string(*tagParam))
			if err != nil {
				return artifact.GetClientSetupDetails404JSONResponse{
					NotFoundJSONResponse: artifact.NotFoundJSONResponse(
						*GetErrorResponse(http.StatusNotFound, "tag doesn't exist"),
					),
				}, err
			}
		}
	}

	packageType := string(reg.PackageType)

	return artifact.GetClientSetupDetails200JSONResponse{
		ClientSetupDetailsResponseJSONResponse: *c.GenerateClientSetupDetails(
			ctx, packageType, imageParam, tagParam, regInfo.RegistryRef,
		),
	}, nil
}

func (c *APIController) GenerateClientSetupDetails(
	ctx context.Context,
	packageType string,
	image *artifact.ArtifactParam,
	tag *artifact.VersionParam,
	registryRef string,
) *artifact.ClientSetupDetailsResponseJSONResponse {
	session, _ := request.AuthSessionFrom(ctx)
	username := session.Principal.Email
	loginUsernameLabel := "Username: <USERNAME>"
	loginUsernameValue := "<USERNAME>"
	loginPasswordLabel := "Password: *see step 2*"
	blankString := ""
	if packageType == string(artifact.PackageTypeMAVEN) {
		return c.generateMavenClientSetupDetail(ctx, image, tag, registryRef, username)
	} else if packageType == string(artifact.PackageTypeHELM) {
		header1 := "Login to Helm"
		section1step1Header := "Run this Helm command in your terminal to authenticate the client."
		helmLoginValue := "helm registry login <LOGIN_HOSTNAME>"
		section1step1Commands := []artifact.ClientSetupStepCommand{
			{Label: &blankString, Value: &helmLoginValue},
			{Label: &loginUsernameLabel, Value: &loginUsernameValue},
			{Label: &loginPasswordLabel, Value: &blankString},
		}
		section1step1Type := artifact.ClientSetupStepTypeStatic
		section1step2Header := "For the Password field above, generate an identity token"
		section1step2Type := artifact.ClientSetupStepTypeGenerateToken
		section1Steps := []artifact.ClientSetupStep{
			{
				Header:   &section1step1Header,
				Commands: &section1step1Commands,
				Type:     &section1step1Type,
			},
			{
				Header: &section1step2Header,
				Type:   &section1step2Type,
			},
		}
		section1 := artifact.ClientSetupSection{
			Header: &header1,
		}
		_ = section1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
			Steps: &section1Steps,
		})

		header2 := "Push a version"
		section2step1Header := "Run this Helm push command in your terminal to push a chart in OCI form." +
			" Note: Make sure you add oci:// prefix to the repository URL."
		helmPushValue := "helm push <CHART_TGZ_FILE> oci://<HOSTNAME>/<REGISTRY_NAME>"
		section2step1Commands := []artifact.ClientSetupStepCommand{
			{Label: &blankString, Value: &helmPushValue},
		}
		section2step1Type := artifact.ClientSetupStepTypeStatic
		section2Steps := []artifact.ClientSetupStep{
			{
				Header:   &section2step1Header,
				Commands: &section2step1Commands,
				Type:     &section2step1Type,
			},
		}
		section2 := artifact.ClientSetupSection{
			Header: &header2,
		}
		_ = section2.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
			Steps: &section2Steps,
		})

		header3 := "Pull a version"
		section3step1Header := "Run this Helm command in your terminal to pull a specific chart version."
		helmPullValue := "helm pull oci://<HOSTNAME>/<REGISTRY_NAME>/<IMAGE_NAME> --version <TAG>"
		section3step1Commands := []artifact.ClientSetupStepCommand{
			{Label: &blankString, Value: &helmPullValue},
		}
		section3step1Type := artifact.ClientSetupStepTypeStatic
		section3Steps := []artifact.ClientSetupStep{
			{
				Header:   &section3step1Header,
				Commands: &section3step1Commands,
				Type:     &section3step1Type,
			},
		}
		section3 := artifact.ClientSetupSection{
			Header: &header3,
		}
		_ = section3.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
			Steps: &section3Steps,
		})
		clientSetupDetails := artifact.ClientSetupDetails{
			MainHeader: "Helm Client Setup",
			SecHeader:  "Follow these instructions to install/use Helm artifacts or compatible packages.",
			Sections: []artifact.ClientSetupSection{
				section1,
				section2,
				section3,
			},
		}

		c.replacePlaceholders(ctx, &clientSetupDetails.Sections, username, registryRef, image, tag, "", "")

		return &artifact.ClientSetupDetailsResponseJSONResponse{
			Data:   clientSetupDetails,
			Status: artifact.StatusSUCCESS,
		}
	}
	header1 := "Login to Docker"
	section1step1Header := "Run this Docker command in your terminal to authenticate the client."
	dockerLoginValue := "docker login <LOGIN_HOSTNAME>"
	section1step1Commands := []artifact.ClientSetupStepCommand{
		{Label: &blankString, Value: &dockerLoginValue},
		{Label: &loginUsernameLabel, Value: &loginUsernameValue},
		{Label: &loginPasswordLabel, Value: &blankString},
	}
	section1step1Type := artifact.ClientSetupStepTypeStatic
	section1step2Header := "For the Password field above, generate an identity token"
	section1step2Type := artifact.ClientSetupStepTypeGenerateToken
	section1Steps := []artifact.ClientSetupStep{
		{
			Header:   &section1step1Header,
			Commands: &section1step1Commands,
			Type:     &section1step1Type,
		},
		{
			Header: &section1step2Header,
			Type:   &section1step2Type,
		},
	}
	section1 := artifact.ClientSetupSection{
		Header: &header1,
	}
	_ = section1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &section1Steps,
	})
	header2 := "Pull an image"
	section2step1Header := "Run this Docker command in your terminal to pull image."
	dockerPullValue := "docker pull <HOSTNAME>/<REGISTRY_NAME>/<IMAGE_NAME>:<TAG>"
	section2step1Commands := []artifact.ClientSetupStepCommand{
		{Label: &blankString, Value: &dockerPullValue},
	}
	section2step1Type := artifact.ClientSetupStepTypeStatic
	section2Steps := []artifact.ClientSetupStep{
		{
			Header:   &section2step1Header,
			Commands: &section2step1Commands,
			Type:     &section2step1Type,
		},
	}
	section2 := artifact.ClientSetupSection{
		Header: &header2,
	}
	_ = section2.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &section2Steps,
	})
	header3 := "Retag and Push the image"
	section3step1Header := "Run this Docker command in your terminal to tag the image."
	dockerTagValue := "docker tag <IMAGE_NAME>:<TAG> <HOSTNAME>/<REGISTRY_NAME>/<IMAGE_NAME>:<TAG>"
	section3step1Commands := []artifact.ClientSetupStepCommand{
		{Label: &blankString, Value: &dockerTagValue},
	}
	section3step1Type := artifact.ClientSetupStepTypeStatic
	section3step2Header := "Run this Docker command in your terminal to push the image."
	dockerPushValue := "docker push <HOSTNAME>/<REGISTRY_NAME>/<IMAGE_NAME>:<TAG>"
	section3step2Commands := []artifact.ClientSetupStepCommand{
		{Label: &blankString, Value: &dockerPushValue},
	}
	section3step2Type := artifact.ClientSetupStepTypeStatic
	section3Steps := []artifact.ClientSetupStep{
		{
			Header:   &section3step1Header,
			Commands: &section3step1Commands,
			Type:     &section3step1Type,
		},
		{
			Header:   &section3step2Header,
			Commands: &section3step2Commands,
			Type:     &section3step2Type,
		},
	}
	section3 := artifact.ClientSetupSection{
		Header: &header3,
	}
	_ = section3.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &section3Steps,
	})
	clientSetupDetails := artifact.ClientSetupDetails{
		MainHeader: "Docker Client Setup",
		SecHeader:  "Follow these instructions to install/use Docker artifacts or compatible packages.",
		Sections: []artifact.ClientSetupSection{
			section1,
			section2,
			section3,
		},
	}

	c.replacePlaceholders(ctx, &clientSetupDetails.Sections, username, registryRef, image, tag, "", "")

	return &artifact.ClientSetupDetailsResponseJSONResponse{
		Data:   clientSetupDetails,
		Status: artifact.StatusSUCCESS,
	}
}

func (c *APIController) generateMavenClientSetupDetail(
	ctx context.Context,
	artifactName *artifact.ArtifactParam,
	version *artifact.VersionParam,
	registryRef string,
	username string,
) *artifact.ClientSetupDetailsResponseJSONResponse {
	staticStepType := artifact.ClientSetupStepTypeStatic
	generateTokenStepType := artifact.ClientSetupStepTypeGenerateToken

	section1 := artifact.ClientSetupSection{
		Header:    stringPtr("1. Generate Identity Token"),
		SecHeader: stringPtr("An identity token will serve as the password for uploading and downloading artifacts."),
	}
	_ = section1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("Generate an identity token"),
				Type:   &generateTokenStepType,
			},
		},
	})

	mavenSection1 := artifact.ClientSetupSection{
		Header:    stringPtr("2. Pull a Maven Package"),
		SecHeader: stringPtr("Set default repository in your pom.xml file."),
	}
	_ = mavenSection1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("To set default registry in your pom.xml file by adding the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("<repositories>\n  <repository>\n    <id>maven-dev</id>\n    <url><REGISTRY_URL>/<REGISTRY_NAME></url>\n    <releases>\n      <enabled>true</enabled>\n      <updatePolicy>always</updatePolicy>\n    </releases>\n    <snapshots>\n      <enabled>true</enabled>\n      <updatePolicy>always</updatePolicy>\n    </snapshots>\n  </repository>\n</repositories>"),
					},
				},
			},
			{
				//nolint:lll
				Header: stringPtr("Copy the following your ~/ .m2/setting.xml file for MacOs, or $USERPROFILE$\\ .m2\\settings.xml for Windows to authenticate with token to pull from your Maven registry."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("<settings>\n  <servers>\n    <server>\n      <id>maven-dev</id>\n      <username><USERNAME></username>\n      <password>identity-token</password>\n    </server>\n  </servers>\n</settings>"),
					},
				},
			},
			{
				//nolint:lll
				Header: stringPtr("Add a dependency to the project's pom.xml (replace <GROUP_ID>, <ARTIFACT_ID> & <VERSION> with your own):"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("<dependency>\n  <groupId><GROUP_ID></groupId>\n  <artifactId><ARTIFACT_ID></artifactId>\n  <version><VERSION></version>\n</dependency>"),
					},
				},
			},
			{
				Header: stringPtr("Install dependencies in pom.xml file"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("mvn install"),
					},
				},
			},
		},
	})

	mavenSection2 := artifact.ClientSetupSection{
		Header:    stringPtr("3. Push a Maven Package"),
		SecHeader: stringPtr("Set default repository in your pom.xml file."),
	}

	_ = mavenSection2.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("To set default registry in your pom.xml file by adding the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("<distributionManagement>\n  <snapshotRepository>\n    <id>maven-dev</id>\n    <url><REGISTRY_URL>/<REGISTRY_NAME></url>\n  </snapshotRepository>\n  <repository>\n    <id>maven-dev</id>\n    <url><REGISTRY_URL>/<REGISTRY_NAME></url>\n  </repository>\n</distributionManagement>"),
					},
				},
			},
			{
				//nolint:lll
				Header: stringPtr("Copy the following your ~/ .m2/setting.xml file for MacOs, or $USERPROFILE$\\ .m2\\settings.xml for Windows to authenticate with token to push to your Maven registry."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("<settings>\n  <servers>\n    <server>\n      <id>maven-dev</id>\n      <username><USERNAME></username>\n      <password>identity-token</password>\n    </server>\n  </servers>\n</settings>"),
					},
				},
			},
			{
				Header: stringPtr("Publish package to your Maven registry."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("mvn deploy"),
					},
				},
			},
		},
	})

	gradleSection1 := artifact.ClientSetupSection{
		Header:    stringPtr("2. Pull a Gradle Package"),
		SecHeader: stringPtr("Set default repository in your build.gradle file."),
	}
	_ = gradleSection1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("Set the default registry in your project’s build.gradle by adding the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("repositories{\n    maven{\n      url “<REGISTRY_URL>/<REGISTRY_NAME>”\n\n      credentials {\n         username “<USERNAME>”\n         password “identity-token”\n      }\n   }\n}"),
					},
				},
			},
			{
				//nolint:lll
				Header: stringPtr("As this is a private registry, you’ll need to authenticate. Create or add to the ~/.gradle/gradle.properties file with the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("repositoryUser=<USERNAME>\nrepositoryPassword={{identity-token}}"),
					},
				},
			},
			{
				Header: stringPtr("Add a dependency to the project’s build.gradle"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("dependencies {\n  implementation ‘<GROUP_ID>:<ARTIFACT_ID>:<VERSION>’\n}"),
					},
				},
			},
			{
				Header: stringPtr("Install dependencies in build.gradle file"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("gradlew build     // Linux or OSX\n gradlew.bat build  // Windows"),
					},
				},
			},
		},
	})

	gradleSection2 := artifact.ClientSetupSection{
		Header:    stringPtr("3. Push a Gradle Package"),
		SecHeader: stringPtr("Set default repository in your build.gradle file."),
	}

	_ = gradleSection2.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("Add a maven publish plugin configuration to the project’s build.gradle."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("publishing {\n    publications {\n        maven(MavenPublication) {\n            groupId = '<GROUP_ID>'\n            artifactId = '<ARTIFACT_ID>'\n            version = '<VERSION>'\n\n            from components.java\n        }\n    }\n}"),
					},
				},
			},
			{
				Header: stringPtr("Publish package to your Maven registry."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("gradlew publish"),
					},
				},
			},
		},
	})

	sbtSection1 := artifact.ClientSetupSection{
		Header:    stringPtr("2. Pull a Sbt/Scala Package"),
		SecHeader: stringPtr("Set default repository in your build.sbt file."),
	}
	_ = sbtSection1.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("Set the default registry in your project’s build.sbt by adding the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("resolver += “Harness Registry” at “<REGISTRY_URL>/<REGISTRY_NAME>”\ncredentials += Credentials(Path.userHome / “.sbt” / “.Credentials”)"),
					},
				},
			},
			{
				//nolint:lll
				Header: stringPtr("As this is a private registry, you’ll need to authenticate. Create or add to the ~/.sbt/.credentials file with the following:"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						//nolint:lll
						Value: stringPtr("realm=Harness Registry\nhost=<LOGIN_HOSTNAME>\nuser=<USERNAME>\npassword={{identity-token}}"),
					},
				},
			},
			{
				Header: stringPtr("Add a dependency to the project’s build.sbt"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("libraryDependencies += “<GROUP_ID>” % “<ARTIFACT_ID>” % “<VERSION>”"),
					},
				},
			},
			{
				Header: stringPtr("Install dependencies in build.sbt file"),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("sbt update"),
					},
				},
			},
		},
	})

	sbtSection2 := artifact.ClientSetupSection{
		Header:    stringPtr("3. Push a Sbt/Scala Package"),
		SecHeader: stringPtr("Set default repository in your build.sbt file."),
	}

	_ = sbtSection2.FromClientSetupStepConfig(artifact.ClientSetupStepConfig{
		Steps: &[]artifact.ClientSetupStep{
			{
				Header: stringPtr("Add publish configuration to the project’s build.sbt."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("publishTo := Some(\"Harness Registry\" at \"<REGISTRY_URL>/<REGISTRY_NAME>\")"),
					},
				},
			},
			{
				Header: stringPtr("Publish package to your Maven registry."),
				Type:   &staticStepType,
				Commands: &[]artifact.ClientSetupStepCommand{
					{
						Value: stringPtr("sbt publish"),
					},
				},
			},
		},
	})

	section2 := artifact.ClientSetupSection{}
	_ = section2.FromTabSetupStepConfig(artifact.TabSetupStepConfig{
		Tabs: &[]artifact.TabSetupStep{
			{
				Header: stringPtr("Maven"),
				Sections: &[]artifact.ClientSetupSection{
					mavenSection1,
					mavenSection2,
				},
			},
			{
				Header: stringPtr("Gradle"),
				Sections: &[]artifact.ClientSetupSection{
					gradleSection1,
					gradleSection2,
				},
			},
			{
				Header: stringPtr("Sbt/Scala"),
				Sections: &[]artifact.ClientSetupSection{
					sbtSection1,
					sbtSection2,
				},
			},
		},
	})

	clientSetupDetails := artifact.ClientSetupDetails{
		MainHeader: "Maven Client Setup",
		SecHeader:  "Follow these instructions to install/use Maven artifacts or compatible packages.",
		Sections: []artifact.ClientSetupSection{
			section1,
			section2,
		},
	}
	groupID := ""
	if artifactName != nil {
		parts := strings.Split(string(*artifactName), ":")
		if len(parts) == 2 {
			groupID = parts[0]
			*artifactName = artifact.ArtifactParam(parts[1])
		}
	}

	rootSpace, _, _ := paths.DisectRoot(registryRef)
	registryURL := c.URLProvider.RegistryURL(ctx, "maven", rootSpace)

	//nolint:lll
	c.replacePlaceholders(ctx, &clientSetupDetails.Sections, username, registryRef, artifactName, version, registryURL, groupID)

	return &artifact.ClientSetupDetailsResponseJSONResponse{
		Data:   clientSetupDetails,
		Status: artifact.StatusSUCCESS,
	}
}

func (c *APIController) replacePlaceholders(
	ctx context.Context,
	clientSetupSections *[]artifact.ClientSetupSection,
	username string,
	regRef string,
	image *artifact.ArtifactParam,
	tag *artifact.VersionParam,
	registryURL string,
	groupID string,
) {
	for i := range *clientSetupSections {
		tab, err := (*clientSetupSections)[i].AsTabSetupStepConfig()
		if err != nil || tab.Tabs == nil {
			c.replacePlaceholdersInSection(ctx, &(*clientSetupSections)[i], username, regRef, image, tag, groupID, registryURL)
		} else {
			for j := range *tab.Tabs {
				c.replacePlaceholders(ctx, (*tab.Tabs)[j].Sections, username, regRef, image, tag, groupID, registryURL)
			}
			_ = (*clientSetupSections)[i].FromTabSetupStepConfig(tab)
		}
	}
}

func (c *APIController) replacePlaceholdersInSection(
	ctx context.Context,
	clientSetupSection *artifact.ClientSetupSection,
	username string,
	regRef string,
	image *artifact.ArtifactParam,
	tag *artifact.VersionParam,
	registryURL string,
	groupID string,
) {
	rootSpace, _, _ := paths.DisectRoot(regRef)
	_, registryName, _ := paths.DisectLeaf(regRef)
	hostname := common.TrimURLScheme(c.URLProvider.RegistryURL(ctx, rootSpace))

	sec, err := clientSetupSection.AsClientSetupStepConfig()
	if err != nil || sec.Steps == nil {
		return
	}
	for _, st := range *sec.Steps {
		if st.Commands == nil {
			continue
		}
		for j := range *st.Commands {
			replaceText(username, st, j, hostname, registryName, image, tag, registryURL, groupID)
		}
	}
	_ = clientSetupSection.FromClientSetupStepConfig(sec)
}

func replaceText(
	username string,
	st artifact.ClientSetupStep,
	i int,
	hostname string,
	repoName string,
	image *artifact.ArtifactParam,
	tag *artifact.VersionParam,
	registryURL string,
	groupID string,
) {
	if username != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<USERNAME>", username))
		if (*st.Commands)[i].Label != nil {
			(*st.Commands)[i].Label = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Label, "<USERNAME>", username))
		}
	}
	if groupID != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<GROUP_ID>", groupID))
	}
	if registryURL != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<REGISTRY_URL>", registryURL))
	}
	if hostname != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<HOSTNAME>", hostname))
	}
	if hostname != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value,
			"<LOGIN_HOSTNAME>", common.GetHost(hostname)))
	}
	if repoName != "" {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<REGISTRY_NAME>", repoName))
	}
	if image != nil {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<IMAGE_NAME>", string(*image)))
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<ARTIFACT_ID>", string(*image)))
	}
	if tag != nil {
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<TAG>", string(*tag)))
		(*st.Commands)[i].Value = stringPtr(strings.ReplaceAll(*(*st.Commands)[i].Value, "<VERSION>", string(*tag)))
	}
}

func stringPtr(s string) *string {
	return &s
}

group "default" {
  targets = ["tidsapparat"]
}

target "tidsapparat" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = ["linux/amd64", "linux/arm64"]

  labels = {
    "org.opencontainers.image.url" = "https://github.com/potibm/tidsapparat"
    "org.opencontainers.image.source" = "https://github.com/potibm/tidsapparat"
    "org.opencontainers.image.documentation" = "https://github.com/potibm/tidsapparat/tree/main/doc"
    "org.opencontainers.image.authors" = "potibm"
  }
  
  annotations = [
    "index,manifest:org.opencontainers.image.title=Tidsapparat",
    "index,manifest:org.opencontainers.image.description=A party information system for the bigscreen at demoparties.",
    "index,manifest:org.opencontainers.image.url=https://github.com/potibm/tidsapparat",
    "index,manifest:org.opencontainers.image.source=https://github.com/potibm/tidsapparat",
    "index,manifest:org.opencontainers.image.documentation=https://github.com/potibm/tidsapparat/tree/main/doc",
    "index,manifest:org.opencontainers.image.licenses=MIT",
    "index,manifest:org.opencontainers.image.authors=potibm"
  ]
}

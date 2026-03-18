plugins {
    id("java")
    id("org.jetbrains.kotlin.jvm") version "2.2.0"
    id("org.jetbrains.intellij") version "1.17.2"
}

group = "com.ink.yht.sprout"
version = "1.0.0"

repositories {
    mavenCentral()
    maven(url = "https://www.jetbrains.com/intellij-repository/releases")
    maven(url = "https://cache-redirector.jetbrains.com/intellij-dependencies")
}

intellij {
    type.set("IC")
    version.set("2024.2")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
}

sourceSets {
    main {
        java {
            srcDirs("src/main/kotlin", "gen")
        }
    }
}

tasks {
    withType<JavaCompile> {
        options.release.set(17)
    }
    withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile> {
        compilerOptions {
            jvmTarget.set(org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_17)
        }
    }
    patchPluginXml {
        sinceBuild.set("242")
        untilBuild.set("253.*")
    }
}

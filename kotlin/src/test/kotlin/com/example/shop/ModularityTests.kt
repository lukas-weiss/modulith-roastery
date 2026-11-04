package com.example.shop

import org.junit.jupiter.api.Test
import org.springframework.modulith.core.ApplicationModules
import org.springframework.modulith.docs.Documenter

class ModularityTests {

    private val modules = ApplicationModules.of(ShopApplication::class.java)

    /** Fails if a module uses another module's internals or if modules depend on each other in a cycle. */
    @Test
    fun `module structure is valid`() {
        modules.forEach(::println)
        modules.verify()
    }

    /** Writes PlantUML diagrams and module canvases to build/spring-modulith-docs. */
    @Test
    fun `write documentation`() {
        Documenter(modules)
            .writeModulesAsPlantUml()
            .writeIndividualModulesAsPlantUml()
            .writeModuleCanvases()
    }
}

import * as THREE from "three"
import { OrbitControls } from "three/addons/controls/OrbitControls.js"
import { treads, risers } from "./stair.mjs"
import { EffectComposer } from "three/addons/postprocessing/EffectComposer.js"
import { RenderPass } from "three/addons/postprocessing/RenderPass.js"
import { GlitchPass } from "three/addons/postprocessing/GlitchPass.js"
import { OutputPass } from "three/addons/postprocessing/OutputPass.js"
import { RenderPixelatedPass } from "three/addons/postprocessing/RenderPixelatedPass.js"
import { UnrealBloomPass } from 'three/addons/postprocessing/UnrealBloomPass.js'

const scene = new THREE.Scene()
const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000)
const renderer = new THREE.WebGLRenderer()
renderer.setSize(window.innerWidth, window.innerHeight)
document.body.appendChild(renderer.domElement)

const controls = new OrbitControls(camera, renderer.domElement)

controls.enableRotate = true

camera.position.set(0, -5, 2.5)
camera.lookAt(0, 15 * 0.255, 0)

controls.update()




const axesHelper = new THREE.AxesHelper(15)
scene.add(axesHelper)

const light = new THREE.AmbientLight(0x404040)
scene.add(light)


let composer = new EffectComposer(renderer)
const pixpass = new RenderPixelatedPass(3, scene, camera)
composer.addPass(pixpass)
const glitchpass = new GlitchPass()
//composer.addPass(glitchpass)
//const upass = new UnrealBloomPass(new THREE.Vector2(window.innerWidth, window.innerHeight), 1.5, 0.4, 0.85)
//upass.threshold = 0
//upass.strength = 0.5
//upass.radius = 0.5
//composer.addPass(upass)
const outputpass = new OutputPass()
composer.addPass(outputpass)

let steps = 15
let going = 255
let rise = 190
let width = 1000

let ts = treads(steps, going, rise, 21, width)
let rs = risers(steps, going, rise, 15, width)
scene.add(ts)
scene.add(rs)

document.getElementById("steps").addEventListener("change", function(event) {
	console.log(event.target.value)
	steps = event.target.value
	scene.remove(ts)
	scene.remove(rs)
	ts = treads(steps, going, rise, 21, width)
	rs = risers(steps, going, rise, 15, width)
	scene.add(ts)
	scene.add(rs)
})

// render loop
function animate() {
	requestAnimationFrame(animate)
	composer.render(scene, camera)
}
renderer.setAnimationLoop(animate)

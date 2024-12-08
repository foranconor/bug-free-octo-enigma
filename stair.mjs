import { BoxGeometry, MeshBasicMaterial, Mesh, Group } from "three";

const SCALE = 0.01
const NOSING = 25

export function treads(n, g, r, t, w) {
	const material = new MeshBasicMaterial({
		color: 0x049ef4,
	})
	const group = new Group()
	for (let i = 0; i < n; i++) {
		const geo = new BoxGeometry(w * SCALE, t * SCALE, (g + NOSING) * SCALE)
		const mes = new Mesh(geo, material)
		mes.position.set(
			0,
			r * SCALE * (i + 1) + t * SCALE / 2,
			g * SCALE * i + (g + NOSING) * SCALE / 2,
		)
		group.add(mes)
	}
	return group
}

export function risers(n, g, r, t, w) {
	const material = new MeshBasicMaterial({
		color: 0xa412f3,
	})
	const group = new Group()
	for (let i = 0; i < n; i++) {
		const geo = new BoxGeometry(w * SCALE, r * SCALE, t * SCALE)
		const mes = new Mesh(geo, material)
		mes.position.set(
			0,
			r * SCALE * i + r * SCALE / 2,
			g * SCALE * i + t * SCALE / 2 + NOSING * SCALE,
		)
		group.add(mes)
	}
	return group
}

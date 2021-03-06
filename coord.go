// Public domain.

// Package coord implements 3D cartesian and 2D equatorial coordinates.
//
// The method sets are not comprehensive, but include what has been
// needed for digest2.
package coord

import (
	"math"

	"github.com/soniakeys/unit"
)

const twoPi = 2 * math.Pi

// Cart represents a general purpose 3D cartesian coordinate.
type Cart struct {
	X, Y, Z float64
}

// Neg sets z = -a, returns z.
func (z *Cart) Neg(a *Cart) *Cart {
	z.X = -a.X
	z.Y = -a.Y
	z.Z = -a.Z
	return z
}

// Add sets z = a1 + a2, returns z.
func (z *Cart) Add(a1, a2 *Cart) *Cart {
	z.X = a1.X + a2.X
	z.Y = a1.Y + a2.Y
	z.Z = a1.Z + a2.Z
	return z
}

// Sub sets z = a1 - a2, returns z.
func (z *Cart) Sub(a1, a2 *Cart) *Cart {
	z.X = a1.X - a2.X
	z.Y = a1.Y - a2.Y
	z.Z = a1.Z - a2.Z
	return z
}

// MulScalar performs element-wise z = a * scalar b, returns z.
func (z *Cart) MulScalar(a *Cart, b float64) *Cart {
	z.X = a.X * b
	z.Y = a.Y * b
	z.Z = a.Z * b
	return z
}

// RotateX rotates the coordinate system around the X axis using the sine and
// cosine of a rotation angle.
//
// (This is useful for translating between equatorial and ecliptic coordinates,
// for example.)
//
// It sets z = a with coordinates rotated by sin, cos; returns z.
func (z *Cart) RotateX(a *Cart, sin, cos float64) *Cart {
	z.X, z.Y, z.Z = a.X, a.Z*sin+a.Y*cos, a.Z*cos-a.Y*sin
	return z
}

// Dot returns the dot product of its argument and the receiver.
func (a1 *Cart) Dot(a2 *Cart) float64 {
	return a1.X*a2.X + a1.Y*a2.Y + a1.Z*a2.Z
}

// Square = a.Dot(a).
func (a *Cart) Square() float64 {
	return a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

// Cross computes the cross product.  It sets z = a × b, returns z.
func (z *Cart) Cross(a, b *Cart) *Cart {
	z.X, z.Y, z.Z =
		a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X
	return z
}

// Equa are equatorial coordinates, referenced to the Earth's rotational axis.
type Equa struct {
	RA  unit.RA    // Right ascension (α)
	Dec unit.Angle // Declination (δ)
}

// Sphr are general purpose spherical coordinates.
type Sphr struct {
	Lon unit.Angle
	Lat unit.Angle
}

// slice types

type CartS []Cart
type EquaS []Equa
type SphrS []Sphr

// FromSphr converts spherical coordinates (Lon, Lat) to a
// cartesian unit vector (X, Y, Z).  The receiver is returned.
func (c *Cart) FromSphr(s *Sphr) *Cart {
	sLon, cLon := s.Lon.Sincos()
	sLat, cLat := s.Lat.Sincos()
	c.X = cLat * cLon
	c.Y = cLat * sLon
	c.Z = sLat
	return c
}

// FromEquaS converts equatorial slice to cartesian slice.
// Receiver length is adjusted to the length of the parameter.
// The receiver is returned.
func (cp *CartS) FromEquaS(e EquaS) CartS {
	c := *cp
	if cap(c) < len(e) {
		c = make(CartS, len(e))
	} else {
		c = c[:len(e)]
	}
	var s Sphr
	for i := range e {
		c[i].FromSphr(s.FromEqua(&e[i]))
	}
	*cp = c
	return c
}

func (s *Sphr) FromEqua(e *Equa) *Sphr {
	s.Lon = e.RA.Angle()
	s.Lat = e.Dec
	return s
}

// This is the copying slice type conversion your mother warned you about.
func (e EquaS) SphrS() SphrS {
	s := make(SphrS, len(e))
	for i := range e {
		s[i].FromEqua(&e[i])
	}
	return s
}

// FromSphrS converts spherical slice to cartesian slice.
// Receiver length is adjusted to the length of the parameter.
// The receiver is returned.
func (cp *CartS) FromSphrS(s SphrS) CartS {
	c := *cp
	if cap(c) < len(s) {
		c = make(CartS, len(s))
	} else {
		c = c[:len(s)]
	}
	for i, s1 := range s {
		c[i].FromSphr(&s1)
	}
	*cp = c
	return c
}

// FromCart cartesian vector to equitorial coordinates.
// The receiver is returned.
func (s *Equa) FromCart(c *Cart) *Equa {
	s.RA = unit.RAFromRad(math.Atan2(c.Y, c.X))
	s.Dec = unit.Angle(math.Asin(c.Z))
	return s
}

// FromCart cartesian vector to spherical coordinates.
// The receiver is returned.
func (s *Sphr) FromCart(c *Cart) *Sphr {
	s.Lon = unit.Angle(math.Atan2(c.Y, c.X))
	s.Lat = unit.Angle(math.Asin(c.Z))
	return s
}

// FromCartS converts cartesian slice to spherical slice.
// Receiver length is adjusted to the length of the parameter.
// The receiver is returned.
func (sp *SphrS) FromCartS(c CartS) SphrS {
	s := *sp
	if cap(s) < len(c) {
		s = make(SphrS, len(c))
	} else {
		s = s[:len(c)]
	}
	for i, c1 := range c {
		s[i].FromCart(&c1)
	}
	*sp = s
	return s
}

// M3 represents a 3 × 3 matrix as a flat array.
type M3 [9]float64

// Mult3 does matrix multiplication.
//
// It sets z = rm × a, returns z.
//
// As a typical use, if rm represents a rotation matrix, then Mult3(rm, a)
// rotates the vector a by rm.
func (z *Cart) Mult3(rm *M3, a *Cart) *Cart {
	*z = Cart{
		rm[0]*a.X + rm[1]*a.Y + rm[2]*a.Z,
		rm[3]*a.X + rm[4]*a.Y + rm[5]*a.Z,
		rm[6]*a.X + rm[7]*a.Y + rm[8]*a.Z,
	}
	return z
}

// Mult3S broadcasts Cart.Mult3 to a slice.
func (zp *CartS) Mult3S(rm *M3, a CartS) CartS {
	z := *zp
	if cap(z) < len(a) {
		z = make(CartS, len(a))
	} else {
		z = z[:len(z)]
	}
	for i, a1 := range a {
		z[i].Mult3(rm, &a1)
	}
	*zp = z
	return z
}

// Transpose will transpose in place or to separate result.
// Sets z = transpose(a), returns z.
func (z *M3) Transpose(a *M3) *M3 {
	z[0], z[1], z[3] = a[0], a[3], a[1]
	z[2], z[4], z[6] = a[6], a[4], a[2]
	z[5], z[7], z[8] = a[7], a[5], a[8]
	return z
}

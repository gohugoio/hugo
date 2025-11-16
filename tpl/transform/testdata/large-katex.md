---
title: Large KaTeX
source:  https://math.stackexchange.com/questions/8337/different-methods-to-compute-sum-limits-k-1-infty-frac1k2-basel-pro
license: https://creativecommons.org/licenses/by-sa/4.0/
---


As I have heard people did not trust Euler when he first discovered the formula (solution of the Basel problem) $$\zeta(2)=\sum_{k=1}^\infty \frac{1}{k^2}=\frac{\pi^2}{6}.$$ However, Euler was Euler and he gave other proofs.
I believe many of you know some nice proofs of this, can you please share it with us?

Hans Lundmark
OK, here's my favorite. I thought of this after reading a proof from the book "Proofs from the book" by Aigner & Ziegler, but later I found more or less the same proof as mine in a paper published a few years earlier by Josef Hofbauer. On Robin's list, the proof most similar to this is number 9 (EDIT: ...which is actually the proof that I read in Aigner & Ziegler).
When $0 < x < \pi/2$ we have $0<\sin x < x < \tan x$ and thus $$\frac{1}{\tan^2 x} < \frac{1}{x^2} < \frac{1}{\sin^2 x}.$$ Note that $1/\tan^2 x = 1/\sin^2 x - 1$. Split the interval $(0,\pi/2)$ into $2^n$ equal parts, and sum the inequality over the (inner) "gridpoints" $x_k=(\pi/2) \cdot (k/2^n)$: $$\sum_{k=1}^{2^n-1} \frac{1}{\sin^2 x_k} - \sum_{k=1}^{2^n-1} 1 < \sum_{k=1}^{2^n-1} \frac{1}{x_k^2} < \sum_{k=1}^{2^n-1} \frac{1}{\sin^2 x_k}.$$ Denoting the sum on the right-hand side by $S_n$, we can write this as $$S_n - (2^n - 1) < \sum_{k=1}^{2^n-1} \left( \frac{2 \cdot 2^n}{\pi} \right)^2 \frac{1}{k^2} < S_n.$$
Although $S_n$ looks like a complicated sum, it can actually be computed fairly easily. To begin with, $$\frac{1}{\sin^2 x} + \frac{1}{\sin^2 (\frac{\pi}{2}-x)} = \frac{\cos^2 x + \sin^2 x}{\cos^2 x \cdot \sin^2 x} = \frac{4}{\sin^2 2x}.$$ Therefore, if we pair up the terms in the sum $S_n$ except the midpoint $\pi/4$ (take the point $x_k$ in the left half of the interval $(0,\pi/2)$ together with the point $\pi/2-x_k$ in the right half) we get 4 times a sum of the same form, but taking twice as big steps so that we only sum over every other gridpoint; that is, over those gridpoints that correspond to splitting the interval into $2^{n-1}$ parts. And the midpoint $\pi/4$ contributes with $1/\sin^2(\pi/4)=2$ to the sum. In short, $$S_n = 4 S_{n-1} + 2.$$ Since $S_1=2$, the solution of this recurrence is $$S_n = \frac{2(4^n-1)}{3}.$$ (For example like this: the particular (constant) solution $(S_p)_n = -2/3$ plus the general solution to the homogeneous equation $(S_h)_n = A \cdot 4^n$, with the constant $A$ determined by the initial condition $S_1=(S_p)_1+(S_h)_1=2$.)
We now have $$ \frac{2(4^n-1)}{3} - (2^n-1) \leq \frac{4^{n+1}}{\pi^2} \sum_{k=1}^{2^n-1} \frac{1}{k^2} \leq \frac{2(4^n-1)}{3}.$$ Multiply by $\pi^2/4^{n+1}$ and let $n\to\infty$. This squeezes the partial sums between two sequences both tending to $\pi^2/6$. Voilà!

Américo Tavares
We can use the function $f(x)=x^{2}$ with $-\pi \leq x\leq \pi $ and find its expansion into a trigonometric Fourier series
$$\dfrac{a_{0}}{2}+\sum_{n=1}^{\infty }(a_{n}\cos nx+b_{n}\sin nx),$$
which is periodic and converges to $f(x)$ in $[-\pi, \pi] $.
Observing that $f(x)$ is even, it is enough to determine the coefficients
$$a_{n}=\dfrac{1}{\pi }\int_{-\pi }^{\pi }f(x)\cos nx\;dx\qquad n=0,1,2,3,...,$$
because
$$b_{n}=\dfrac{1}{\pi }\int_{-\pi }^{\pi }f(x)\sin nx\;dx=0\qquad n=1,2,3,... .$$
For $n=0$ we have
$$a_{0}=\dfrac{1}{\pi }\int_{-\pi }^{\pi }x^{2}dx=\dfrac{2}{\pi }\int_{0}^{\pi }x^{2}dx=\dfrac{2\pi ^{2}}{3}.$$
And for $n=1,2,3,...$ we get
$$a_{n}=\dfrac{1}{\pi }\int_{-\pi }^{\pi }x^{2}\cos nx\;dx$$
$$=\dfrac{2}{\pi }\int_{0}^{\pi }x^{2}\cos nx\;dx=\dfrac{2}{\pi }\times \dfrac{ 2\pi }{n^{2}}(-1)^{n}=(-1)^{n}\dfrac{4}{n^{2}},$$
because
$$\int x^2\cos nx\;dx=\dfrac{2x}{n^{2}}\cos nx+\left( \frac{x^{2}}{ n}-\dfrac{2}{n^{3}}\right) \sin nx.$$
Thus
$$f(x)=\dfrac{\pi ^{2}}{3}+\sum_{n=1}^{\infty }\left( (-1)^{n}\dfrac{4}{n^{2}} \cos nx\right) .$$
Since $f(\pi )=\pi ^{2}$, we obtain
$$\pi ^{2}=\dfrac{\pi ^{2}}{3}+\sum_{n=1}^{\infty }\left( (-1)^{n}\dfrac{4}{ n^{2}}\cos \left( n\pi \right) \right) $$
$$\pi ^{2}=\dfrac{\pi ^{2}}{3}+4\sum_{n=1}^{\infty }\left( (-1)^{n}(-1)^{n} \dfrac{1}{n^{2}}\right) $$
$$\pi ^{2}=\dfrac{\pi ^{2}}{3}+4\sum_{n=1}^{\infty }\dfrac{1}{n^{2}}.$$
Therefore
$$\sum_{n=1}^{\infty }\dfrac{1}{n^{2}}=\dfrac{\pi ^{2}}{4}-\dfrac{\pi ^{2}}{12}= \dfrac{\pi ^{2}}{6}$$
Second method (available on-line a few years ago) by Eric Rowland. From
$$\log (1-t)=-\sum_{n=1}^{\infty}\dfrac{t^n}{n}$$
and making the substitution $t=e^{ix}$ one gets the series expansion
$$w=\text{Log}(1-e^{ix})=-\sum_{n=1}^{\infty }\dfrac{e^{inx}}{n}=-\sum_{n=1}^{ \infty }\dfrac{1}{n}\cos nx-i\sum_{n=1}^{\infty }\dfrac{1}{n}\sin nx,$$
whose radius of convergence is $1$. Now if we take the imaginary part of both sides, the RHS becomes
$$\Im w=-\sum_{n=1}^{\infty }\dfrac{1}{n}\sin nx,$$
and the LHS
$$\Im w=\arg \left( 1-\cos x-i\sin x\right) =\arctan \dfrac{-\sin x}{ 1-\cos x}.$$
Since
$$\arctan \dfrac{-\sin x}{1-\cos x}=-\arctan \dfrac{2\sin \dfrac{x}{2}\cdot \cos \dfrac{x}{2}}{2\sin ^{2}\dfrac{x}{2}}$$
$$=-\arctan \cot \dfrac{x}{2}=-\arctan \tan \left( \dfrac{\pi }{2}-\dfrac{x}{2} \right) =\dfrac{x}{2}-\dfrac{\pi }{2},$$
the following expansion holds
$$\dfrac{\pi }{2}-\frac{x}{2}=\sum_{n=1}^{\infty }\dfrac{1}{n}\sin nx.\qquad (\ast )$$
Integrating the identity $(\ast )$, we obtain
$$\dfrac{\pi }{2}x-\dfrac{x^{2}}{4}+C=-\sum_{n=1}^{\infty }\dfrac{1}{n^{2}}\cos nx.\qquad (\ast \ast )$$
Setting $x=0$, we get the relation between $C$ and $\zeta (2)$
$$C=-\sum_{n=1}^{\infty }\dfrac{1}{n^{2}}=-\zeta (2).$$
And for $x=\pi $, since
$$\zeta (2)=2\sum_{n=1}^{\infty }\dfrac{(-1)^{n-1}}{n^{2}},$$
we deduce
$$\dfrac{\pi ^{2}}{4}+C=-\sum_{n=1}^{\infty }\dfrac{1}{n^{2}}\cos n\pi =\sum_{n=1}^{\infty }\dfrac{(-1)^{n-1}}{n^{2}}=\dfrac{1}{2}\zeta (2)=-\dfrac{1}{ 2}C.$$
Solving for $C$
$$C=-\dfrac{\pi ^{2}}{6},$$
we thus prove
$$\zeta (2)=\dfrac{\pi ^{2}}{6}.$$
Note: this 2nd method can generate all the zeta values $\zeta (2n)$ by integrating repeatedly $(\ast\ast )$. This is the reason why I appreciate it. Unfortunately it does not work for $\zeta (2n+1)$.
Note also the $$C=-\dfrac{\pi ^{2}}{6}$$ can be obtained by integrating $(\ast\ast )$ and substitute $$x=0,x=\pi$$ respectively.

AD.
Here is an other one which is more or less what Euler did in one of his proofs.
The function $\sin x$ where $x\in\mathbb{R}$ is zero exactly at $x=n\pi$ for each integer $n$. If we factorized it as an infinite product we get
$$\sin x = \cdots\left(1+\frac{x}{3\pi}\right)\left(1+\frac{x}{2\pi}\right)\left(1+\frac{x}{\pi}\right)x\left(1-\frac{x}{\pi}\right)\left(1-\frac{x}{2\pi}\right)\left(1-\frac{x}{3\pi}\right)\cdots =$$ $$= x\left(1-\frac{x^2}{\pi^2}\right)\left(1-\frac{x^2}{2^2\pi^2}\right)\left(1-\frac{x^2}{3^2\pi^2}\right)\cdots\quad.$$
We can also represent $\sin x$ as a Taylor series at $x=0$:
$$\sin x = x - \frac{x^3}{3!}+\frac{x^5}{5!}-\frac{x^7}{7!}+\cdots\quad.$$
Multiplying the product and identifying the coefficient of $x^3$ we see that
$$\frac{x^3}{3!}=x\left(\frac{x^2}{\pi^2} + \frac{x^2}{2^2\pi^2}+ \frac{x^2}{3^2\pi^2}+\cdots\right)=x^3\sum_{n=1}^{\infty}\frac{1}{n^2\pi^2}$$ or $$\sum_{n=1}^\infty\frac{1}{n^2}=\frac{\pi^2}{6}.$$

Alfredo Z.
Define the following series for $ x > 0 $
$$\frac{\sin x}{x} = 1 - \frac{x^2}{3!}+\frac{x^4}{5!}-\frac{x^6}{7!}+\cdots\quad.$$
Now substitute $ x = \sqrt{y}\ $ to arrive at
$$\frac{\sin \sqrt{y}\ }{\sqrt{y}\ } = 1 - \frac{y}{3!}+\frac{y^2}{5!}-\frac{y^3}{7!}+\cdots\quad.$$
if we find the roots of $\frac{\sin \sqrt{y}\ }{\sqrt{y}\ } = 0 $ we find that
$ y = n^2\pi^2\ $ for $ n \neq 0 $ and $ n $ in the integers
With all of this in mind, recall that for a polynomial
$ P(x) = a_{n}x^n + a_{n-1}x^{n-1} +\cdots+a_{1}x + a_{0} $ with roots $ r_{1}, r_{2}, \cdots , r_{n} $
$$\frac{1}{r_{1}} + \frac{1}{r_{2}} + \cdots + \frac{1}{r_{n}} = -\frac{a_{1}}{a_{0}}$$
Treating the above series for $ \frac{\sin \sqrt{y}\ }{\sqrt{y}\ } $ as polynomial we see that
$$\frac{1}{1^2\pi^2} + \frac{1}{2^2\pi^2} + \frac{1}{3^2\pi^2} + \cdots = -\frac{-\frac{1}{3!}}{1}$$
then multiplying both sides by $ \pi^2 $ gives the desired series.
$$\frac{1}{1^2} + \frac{1}{2^2} + \frac{1}{3^2} + \cdots = \frac{\pi^2}{6}$$

Nameless
This method apparently was used by Tom Apostol in $1983$. I will outline the main ideas of the proof, the details can be found in here or this presentation (page $27$)
Consider
$$\begin{align} \int_{0}^{1} \int_{0}^{1} \frac{1}{1 - xy} dy dx &= \int_{0}^{1} \int_{0}^{1} \sum_{n \geq 0} (xy)^n dy dx \\ &= \sum_{n \geq 0} \int_{0}^{1} \int_{0}^{1} x^n y^n dy dx \\ &= \sum_{n \geq 1} \frac{1}{n^2} \\ \end{align}$$
You can verify that the left hand side is indeed $\frac{\pi^2}{6}$ by letting $x = u - v$ and $y = v + u.$

Qiaochu Yuan
I have two favorite proofs. One is the last proof in Robin Chapman's collection; you really should take a look at it.
The other is a proof that generalizes to the evaluation of $\zeta(2n)$ for all $n$, although I'll do it "Euler-style" to shorten the presentation. The basic idea is that meromorphic functions have infinite partial fraction decompositions that generalize the partial fraction decompositions of rational functions.
The particular function we're interested in is $B(x) = \frac{x}{e^x - 1}$, the exponential generating function of the Bernoulli numbers $B_n$. $B$ is meromorphic with poles at $x = 2 \pi i n, n \in \mathbb{Z}$, and at these poles it has residue $2\pi i n$. It follows that we can write, a la Euler,
$$\frac{x}{e^x - 1} = \sum_{n \in \mathbb{Z}} \frac{2\pi i n}{x - 2 \pi i n} = \sum_{n \in \mathbb{Z}} - \left( \frac{1}{1 - \frac{x}{2\pi i n}} \right).$$
Now we can expand each of the terms on the RHS as a geometric series, again a la Euler, to obtain
$$\frac{x}{e^x - 1} = - \sum_{n \in \mathbb{Z}} \sum_{k \ge 0} \left( \frac{x}{2\pi i n} \right)^k = \sum_{k \ge 0} (-1)^{n+1} \frac{2 \zeta(2n)}{(2\pi )^{2n}} x^{2n}$$
because, after rearranging terms, the sum over odd powers cancels out and the sum over even powers doesn't. (This is one indication of why there is no known closed form for $\zeta(2n+1)$.) Equating terms on both sides, it follows that
$$B_{2n} = (-1)^{n+1} \frac{2 \zeta(2n)}{(2\pi)^{2n}}$$
or
$$\zeta(2n) = (-1)^{n+1} \frac{B_{2n} (2\pi)^{2n}}{2}$$
as desired. To compute $\zeta(2)$ it suffices to compute that $B_2 = \frac{1}{6}$, which then gives the usual answer.
Here is one more nice proof, I learned it from Grisha Mikhalkin:
Lemma: Let $Z$ be a complex curve in $\mathbb{C}^2$. Let $R(Z) \subset \mathbb{R}^2$ be the projection of $Z$ onto its real parts and $I(Z)$ the projection onto its complex parts. If these projections are both one to one, then the area of $R(Z)$ is equal to the area of $I(Z)$.
Proof: There is an obvious map from $R(Z)$ to $I(Z)$, given by lifting $(x_1, x_2) \in R(Z)$ to $(x_1+i y_1, x_2 + i y_2) \in Z$, and then projecting to $(y_1, y_2) \in I(Z)$. We must prove this map has Jacobian $1$. WLOG, translate $(x_1, y_1, x_2, y_2)$ to $(0,0,0,0)$ and let $Z$ obey $\partial z_2/\partial z_1 = a+bi$ near $(0,0)$. To first order, we have $x_2 = a x_1 - b y_1$ and $y_2 = a y_1 + b x_1$. So $y_1 = (a/b) x_1 - (1/b) x_2$ and $y_2 = (a^2 + b^2)/b x_1 - (a/b) x_2$. So the derivative of $(x_1, x_2) \mapsto (y_1, y_2)$ is $\left( \begin{smallmatrix} a/b & - 1/b \\ (a^2 + b^2)/b & -a/b \end{smallmatrix} \right)$ and the Jacobian is $1$. QED
Now, consider the curve $e^{-z_1} + e^{-z_2} = 1$, where $z_1$ and $z_2$ obey the following inequalities: $x_1 \geq 0$, $x_2 \geq 0$, $-\pi \leq y_1 \leq 0$ and $0 \leq y_2 \leq \pi$.
Given a point on $e^{-z_1} + e^{-z_2} = 1$, consider the triangle with vertices at $0$, $e^{-z_1}$ and $e^{-z_1} + e^{-z_2} = 1$. The inequalities on the $y$'s states that the triangle should lie above the real axis; the inequalities on the $x$'s state the horizontal base should be the longest side.
Projecting onto the $x$ coordinates, we see that the triangle exists if and only if the triangle inequality $e^{-x_1} + e^{-x_2} \geq 1$ is obeyed. So $R(Z)$ is the region under the curve $x_2 = - \log(1-e^{-x_1})$. The area under this curve is $$\int_{0}^{\infty} - \log(1-e^{-x}) dx = \int_{0}^{\infty} \sum \frac{e^{-kx}}{k} dx = \sum \frac{1}{k^2}.$$
Now, project onto the $y$ coordinates. Set $(y_1, y_2) = (-\theta_1, \theta_2)$ for convenience, so the angles of the triangle are $(\theta_1, \theta_2, \pi - \theta_1 - \theta_2)$. The largest angle of a triangle is opposite the largest side, so we want $\theta_1$, $\theta_2 \leq \pi - \theta_1 - \theta_2$, plus the obvious inequalities $\theta_1$, $\theta_2 \geq 0$. So $I(Z)$ is the quadrilateral with vertices at $(0,0)$, $(0, \pi/2)$, $(\pi/3, \pi/3)$ and $(\pi/2, 0)$ and, by elementary geometry, this has area $\pi^2/6$.

David Speyer
I'll post the one I know since it is Euler's, and is quite easy and stays in $\mathbb{R}$. (I'm guessing Euler didn't have tools like residues back then).

Peter Tamaroff
Let
$$s = {\sin ^{ - 1}}x$$
Then
$$\int\limits_0^{\frac{\pi }{2}} {sds} = \frac{{{\pi ^2}}}{8}$$
But then
$$\int\limits_0^1 {\frac{{{{\sin }^{ - 1}}x}}{{\sqrt {1 - {x^2}} }}dx} = \frac{{{\pi ^2}}}{8}$$
Since
$${\sin ^{ - 1}}x = \int {\frac{{dx}}{{\sqrt {1 - {x^2}} }}} = x + \frac{1}{2}\frac{{{x^3}}}{3} + \frac{{1 \cdot 3}}{{2 \cdot 4}}\frac{{{x^5}}}{5} + \frac{{1 \cdot 3 \cdot 5}}{{2 \cdot 4 \cdot 6}}\frac{{{x^7}}}{7} + \cdots $$
We have
$$\int\limits_0^1 {\left\{ {\frac{{dx}}{{\sqrt {1 - {x^2}} }}\int {\frac{{dx}}{{\sqrt {1 - {x^2}} }}} } \right\}} = \int\limits_0^1 {\left\{ {x + \frac{1}{2}\frac{{{x^3}}}{3}\frac{{dx}}{{\sqrt {1 - {x^2}} }} + \frac{{1 \cdot 3}}{{2 \cdot 4}}\frac{{{x^5}}}{5}\frac{{dx}}{{\sqrt {1 - {x^2}} }} + \cdots } \right\}} $$
But
$$\int\limits_0^1 {\frac{{{x^{2n + 1}}}}{{\sqrt {1 - {x^2}} }}dx} = \frac{{2n}}{{2n + 1}}\int\limits_0^1 {\frac{{{x^{2n - 1}}}}{{\sqrt {1 - {x^2}} }}dx} $$
which yields
$$\int\limits_0^1 {\frac{{{x^{2n + 1}}}}{{\sqrt {1 - {x^2}} }}dx} = \frac{{\left( {2n} \right)!!}}{{\left( {2n + 1} \right)!!}}$$
since all powers are odd.
This ultimately produces:
$$\frac{{{\pi ^2}}}{8} = 1 + \frac{1}{2}\frac{1}{3}\left( {\frac{2}{3}} \right) + \frac{{1 \cdot 3}}{{2 \cdot 4}}\frac{1}{5}\left( {\frac{{2 \cdot 4}}{{3 \cdot 5}}} \right) + \frac{{1 \cdot 3 \cdot 5}}{{2 \cdot 4 \cdot 6}}\frac{1}{7}\left( {\frac{{2 \cdot 4 \cdot 6}}{{3 \cdot 5 \cdot 7}}} \right) \cdots $$
$$\frac{{{\pi ^2}}}{8} = 1 + \frac{1}{{{3^2}}} + \frac{1}{{{5^2}}} + \frac{1}{{{7^2}}} + \cdots $$
Let
$$1 + \frac{1}{{{2^2}}} + \frac{1}{{{3^2}}} + \frac{1}{{{4^2}}} + \cdots = \omega $$
Then
$$\frac{1}{{{2^2}}} + \frac{1}{{{4^2}}} + \frac{1}{{{6^2}}} + \frac{1}{{{8^2}}} + \cdots = \frac{\omega }{4}$$
Which means
$$\frac{\omega }{4} + \frac{{{\pi ^2}}}{8} = \omega $$
or
$$\omega = \frac{{{\pi ^2}}}{6}$$

Mike Spivey
The most recent issue of The American Mathematical Monthly (August-September 2011, pp. 641-643) has a new proof by Luigi Pace based on elementary probability. Here's the argument.
Let $X_1$ and $X_2$ be independent, identically distributed standard half-Cauchy random variables. Thus their common pdf is $p(x) = \frac{2}{\pi (1+x^2)}$ for $x > 0$.
Let $Y = X_1/X_2$. Then the pdf of $Y$ is, for $y > 0$, $$p_Y(y) = \int_0^{\infty} x p_{X_1} (xy) p_{X_2}(x) dx = \frac{4}{\pi^2} \int_0^\infty \frac{x}{(1+x^2 y^2)(1+x^2)}dx$$ $$=\frac{2}{\pi^2 (y^2-1)} \left[\log \left( \frac{1+x^2 y^2}{1+x^2}\right) \right]_{x=0}^{\infty} = \frac{2}{\pi^2} \frac{\log(y^2)}{y^2-1} = \frac{4}{\pi^2} \frac{\log(y)}{y^2-1}.$$
Since $X_1$ and $X_2$ are equally likely to be the larger of the two, we have $P(Y < 1) = 1/2$. Thus $$\frac{1}{2} = \int_0^1 \frac{4}{\pi^2} \frac{\log(y)}{y^2-1} dy.$$ This is equivalent to $$\frac{\pi^2}{8} = \int_0^1 \frac{-\log(y)}{1-y^2} dy = -\int_0^1 \log(y) (1+y^2+y^4 + \cdots) dy = \sum_{k=0}^\infty \frac{1}{(2k+1)^2},$$ which, as others have pointed out, implies $\zeta(2) = \pi^2/6$.

Hans Lundmark
This is not really an answer, but rather a long comment prompted by David Speyer's answer. The proof that David gives seems to be the one in How to compute $\sum 1/n^2$ by solving triangles by Mikael Passare, although that paper uses a slightly different way of seeing that the area of the region $U_0$ (in Passare's notation) bounded by the positive axes and the curve $e^{-x}+e^{-y}=1$, $$\int_0^{\infty} -\ln(1-e^{-x}) dx,$$ is equal to $\sum_{n\ge 1} \frac{1}{n^2}$.
This brings me to what I really wanted to mention, namely another curious way to see why $U_0$ has that area; I learned this from Johan Wästlund. Consider the region $D_N$ illustrated below for $N=8$:
A shape with area = sum of reciprocal squares
Although it's not immediately obvious, the area of $D_N$ is $\sum_{n=1}^N \frac{1}{n^2}$. Proof: The area of $D_1$ is 1. To get from $D_N$ to $D_{N+1}$ one removes the boxes along the top diagonal, and adds a new leftmost column of rectangles of width $1/(N+1)$ and heights $1/1,1/2,\ldots,1/N$, plus a new bottom row which is the "transpose" of the new column, plus a square of side $1/(N+1)$ in the bottom left corner. The $k$th rectangle from the top in the new column and the $k$th rectangle from the left in the new row (not counting the square) have a combined area which exactly matches the $k$th box in the removed diagonal: $$ \frac{1}{k} \frac{1}{N+1} + \frac{1}{N+1} \frac{1}{N+1-k} = \frac{1}{k} \frac{1}{N+1-k}. $$ Thus the area added in the process is just that of the square, $1/(N+1)^2$. Q.E.D.
(Apparently this shape somehow comes up in connection with the "random assignment problem", where there's an expected value of something which turns out to be $\sum_{n=1}^N \frac{1}{n^2}$.)
Now place $D_N$ in the first quadrant, with the lower left corner at the origin. Letting $N\to\infty$ gives nothing but the region $U_0$: for large $N$ and for $0<\alpha<1$, the upper corner of column number $\lceil \alpha N \rceil$ in $D_N$ lies at $$ (x,y) = \left( \sum_{n=\lceil (1-\alpha) N \rceil}^N \frac{1}{n}, \sum_{n=\lceil \alpha N \rceil}^N \frac{1}{n} \right) \sim \left(\ln\frac{1}{1-\alpha}, \ln\frac{1}{\alpha}\right),$$ hence (in the limit) on the curve $e^{-x}+e^{-y}=1$.

xpaul
Note that $$ \frac{\pi^2}{\sin^2\pi z}=\sum_{n=-\infty}^{\infty}\frac{1}{(z-n)^2} $$ from complex analysis and that both sides are analytic everywhere except $n=0,\pm 1,\pm 2,\cdots$. Then one can obtain $$ \frac{\pi^2}{\sin^2\pi z}-\frac{1}{z^2}=\sum_{n=1}^{\infty}\frac{1}{(z-n)^2}+\sum_{n=1}^{\infty}\frac{1}{(z+n)^2}. $$ Now the right hand side is analytic at $z=0$ and hence $$\lim_{z\to 0}\left(\frac{\pi^2}{\sin^2\pi z}-\frac{1}{z^2}\right)=2\sum_{n=1}^{\infty}\frac{1}{n^2}.$$ Note $$\lim_{z\to 0}\left(\frac{\pi^2}{\sin^2\pi z}-\frac{1}{z^2}\right)=\frac{\pi^2}{3}.$$ Thus $$\sum_{n=1}^{\infty}\frac{1}{n^2}=\frac{\pi^2}{6}.$$

Jack D'Aurizio
Just as a curiosity, a one-line-real-analytic-proof I found by combining different ideas from this thread and this question:
$$\begin{eqnarray*}\zeta(2)&=&\frac{4}{3}\sum_{n=0}^{+\infty}\frac{1}{(2n+1)^2}=\frac{4}{3}\int_{0}^{1}\frac{\log y}{y^2-1}dy\\&=&\frac{2}{3}\int_{0}^{1}\frac{1}{y^2-1}\left[\log\left(\frac{1+x^2 y^2}{1+x^2}\right)\right]_{x=0}^{+\infty}dy\\&=&\frac{4}{3}\int_{0}^{1}\int_{0}^{+\infty}\frac{x}{(1+x^2)(1+x^2 y^2)}dx\,dy\\&=&\frac{4}{3}\int_{0}^{1}\int_{0}^{+\infty}\frac{dx\, dz}{(1+x^2)(1+z^2)}=\frac{4}{3}\cdot\frac{\pi}{4}\cdot\frac{\pi}{2}=\frac{\pi^2}{6}.\end{eqnarray*}$$
Update. By collecting pieces, I have another nice proof. By Euler's acceleration method or just an iterated trick like my $(1)$ here we get: $$ \zeta(2) = \sum_{n\geq 1}\frac{1}{n^2} = \color{red}{\sum_{n\geq 1}\frac{3}{n^2\binom{2n}{n}}}\tag{A}$$ and the last series converges pretty fast. Then we may notice that the last series comes out from a squared arcsine. That just gives another proof of $ \zeta(2)=\frac{\pi^2}{6}$.
A proof of the identity $$\sum_{n\geq 0}\frac{1}{(2n+1)^2}=\frac{\pi}{2}\sum_{k\geq 0}\frac{(-1)^k}{2k+1}=\frac{\pi}{2}\cdot\frac{\pi}{4}$$ is also hidden in tired's answer here. For short, the integral $$ I=\int_{-\infty}^{\infty}e^y\left(\frac{e^y-1}{y^2}-\frac{1}{y}\right)\frac{1}{e^{2y}+1}\,dy $$ is clearly real, so the imaginary part of the sum of residues of the integrand function has to be zero.
Still another way (and a very efficient one) is to exploit the reflection formula for the trigamma function: $$\psi'(1-z)+\psi'(z)=\frac{\pi^2}{\sin^2(\pi z)}$$ immediately leads to: $$\frac{\pi^2}{2}=\psi'\left(\frac{1}{2}\right)=\sum_{n\geq 0}\frac{1}{\left(n+\frac{1}{2}\right)^2}=4\sum_{n\geq 0}\frac{1}{(2n+1)^2}=3\,\zeta(2).$$
2018 update. We may consider that $\mathcal{J}=\int_{0}^{+\infty}\frac{\arctan x}{1+x^2}\,dx = \left[\frac{1}{2}\arctan^2 x\right]_0^{+\infty}=\frac{\pi^2}{8}$.
On the other hand, by Feynman's trick or Fubini's theorem $$ \mathcal{J}=\int_{0}^{+\infty}\int_{0}^{1}\frac{x}{(1+x^2)(1+a^2 x^2)}\,da\,dx = \int_{0}^{1}\frac{-\log a}{1-a^2}\,da $$ and since $\int_{0}^{1}-\log(x)x^n\,dx = \frac{1}{(n+1)^2}$, by expanding $\frac{1}{1-a^2}$ as a geometric series we have $$ \frac{\pi^2}{8}=\mathcal{J}=\sum_{n\geq 0}\frac{1}{(2n+1)^2}. $$

Andrey Rekalo
Here is a complex-analytic proof.
For $z\in D=\mathbb{C}\backslash${$0,1$}, let
$$R(z)=\sum\frac{1}{\log^2 z}$$
where the sum is taken over all branches of the logarithm. Each point in $D$ has a neighbourhood on which the branches of $\log(z)$ are analytic. Since the series converges uniformly away from $z=1$, $R(z)$ is analytic on $D$.
Now a few observations:
(i) Each term of the series tends to $0$ as $z\to0$. Thanks to the uniform convergence this implies that the singularity at $z=0$ is removable and we can set $R(0)=0$.
(ii) The only singularity of $R$ is a double pole at $z=1$ due to the contribution of the principal branch of $\log z$. Moreover, $\lim_{z\to1}(z-1)^2R(z)=1$.
(iii) $R(1/z)=R(z)$.
By (i) and (iii) $R$ is meromorphic on the extended complex plane, therefore it is rational. By (ii) the denominator of $R(z)$ is $(z-1)^2$. Since $R(0)=R(\infty)=0$, the numerator has the form $az$. Then (ii) implies $a=1$, so that $$R(z)=\frac{z}{(z-1)^2}.$$
Now, setting $z=e^{2\pi i w}$ yields $$\sum\limits_{n=-\infty}^{\infty}\frac{1}{(w-n)^2}=\frac{\pi^2}{\sin^2(\pi w)}$$ which implies that $$\sum\limits_{k=0}^{\infty}\frac{1}{(2k+1)^2}=\frac{\pi^2}{8},$$ and the identity $\zeta(2)=\pi^2/6$ follows.
The proof is due to T. Marshall (American Mathematical Monthly, Vol. 117(4), 2010, P. 352).

David Speyer
In response to a request here: Compute $\oint z^{-2k} \cot (\pi z) dz$ where the integral is taken around a square of side $2N+1$. Routine estimates show that the integral goes to $0$ as $N \to \infty$.
Now, let's compute the integral by residues. At $z=0$, the residue is $\pi^{2k-1} q$, where $q$ is some rational number coming from the power series for $\cot$. For example, if $k=1$, then we get $- \pi/3$.
At $m \pi$, for $m \neq 0$, the residue is $z^{-2k} \pi^{-1}$. So $$\pi^{-1} \lim_{N \to \infty} \sum_{-N \leq m \leq N\ m \neq 0} m^{-2k} + \pi^{2k-1} q=0$$ or $$\sum_{m=1}^{\infty} m^{-2k} = -\pi^{2k} q/2$$ as desired. In particular, $\sum m^{-2} = - (\pi^2/3)/2 = \pi^2/6$.
Common variants: We can replace $\cot$ with $\tan$, with $1/(e^{2 \pi i z}-1)$, or with similar formulas.
This is reminiscent of Qiaochu's proof but, rather than actually establishing the relation $\pi^{-1} \cot(\pi z) = \sum (z-n)^{-1}$, one simply establishes that both sides contribute the same residues to a certain integral.

Derek Jennings
Another variation. We make use of the following identity (proved at the bottom of this note):
$$\sum_{k=1}^n \cot^2 \left( \frac {2k-1}{2n} \frac{\pi}{2} \right) = 2n^2 - n. \quad (1)$$
Now $1/\theta > \cot \theta > 1/\theta - \theta/3 > 0$ for $0< \theta< \pi/2 < \sqrt{3}$ and so $$ 1/\theta^2 - 2/3 < \cot^2 \theta < 1/\theta^2. \quad (2)$$
With $\theta_k = (2k-1)\pi/4n,$ summing the inequalities $(2)$ from $k=1$ to $n$ we obtain
$$2n^2 - n < \sum_{k=1}^n \left( \frac{2n}{2k-1}\frac{2}{\pi} \right)^2 < 2n^2 - n + 2n/3.$$
Hence
$$\frac{\pi^2}{16}\frac{2n^2-n}{n^2} < \sum_{k=1}^n \frac{1}{(2k-1)^2} < \frac{\pi^2}{16}\frac{2n^2-n/3}{n^2}.$$
Taking the limit as $n \rightarrow \infty$ we obtain
$$ \sum_{k=1}^\infty \frac{1}{(2k-1)^2} = \frac{\pi^2}{8},$$
from which the result for $\sum_{k=1}^\infty 1/k^2$ follows easily.
To prove $(1)$ we note that
$$ \cos 2n\theta = \text{Re}(\cos\theta + i \sin\theta)^{2n} = \sum_{k=0}^n (-1)^k {2n \choose 2k}\cos^{2n-2k}\theta\sin^{2k}\theta.$$
Therefore
$$\frac{\cos 2n\theta}{\sin^{2n}\theta} = \sum_{k=0}^n (-1)^k {2n \choose 2k}\cot^{2n-2k}\theta.$$
And so setting $x = \cot^2\theta$ we note that
$$f(x) = \sum_{k=0}^n (-1)^k {2n \choose 2k}x^{n-k}$$
has roots $x_j = \cot^2 (2j-1)\pi/4n,$ for $j=1,2,\ldots,n,$ from which $(1)$ follows since ${2n \choose 2n-2} = 2n^2-n.$

xpaul
A short way to get the sum is to use Fourier's expansion of $x^2$ in $x\in(-\pi,\pi)$. Recall that Fourier's expansion of $f(x)$ is $$ \tilde{f}(x)=\frac{1}{2}a_0+\sum_{n=1}^\infty(a_n\cos nx+b_n\sin nx), x\in(-\pi,\pi)$$ where $$ a_0=\frac{2}{\pi}\int_{-\pi}^{\pi}f(x)\;dx, a_n=\frac{2}{\pi}\int_{-\pi}^{\pi}f(x)\cos nx\; dx, b_n=\frac{2}{\pi}\int_{-\pi}^{\pi}f(x)\sin nx\; dx, n=1,2,3,\cdots $$ and $$ \tilde{f}(x)=\frac{f(x-0)+f(x+0)}{2}. $$ Easy calculation shows $$ x^2=\frac{\pi^2}{3}+4\sum_{n=1}^\infty(-1)^n\frac{\cos nx}{n^2}, x\in[-\pi,\pi]. $$ Letting $x=\pi$ in both sides gives $$ \sum_{n=1}^\infty\frac{1}{n^2}=\frac{\pi^2}{6}.$$
Another way to get the sum is to use Parseval's Identity for Fourier's expansion of $x$ in $(-\pi,\pi)$. Recall that Parseval's Identity is $$ \int_{-\pi}^{\pi}|f(x)|^2dx=\frac{1}{2}a_0^2+\sum_{n=1}^\infty(a_n^2+b_n^2). $$ Note $$ x=2\sum_{n=1}^\infty(-1)^{n+1}\frac{\sin nx}{n}, x\in(-\pi,\pi). $$ Using Parseval's Identity gives $$ 4\sum_{n=1}^\infty\frac{1}{n^2}=\int_{-\pi}^{\pi}|x|^2dx$$ or $$ \sum_{n=1}^\infty\frac{1}{n^2}=\frac{\pi^2}{6}.$$

Tomás
I like this one:
Let $f\in Lip(S^{1})$, where $Lip(S^{1})$ is the space of Lipschitz functions on $S^{1}$. So its well defined the number for $k\in \mathbb{Z}$ (called Fourier series of $f$) $$\hat{f}(k)=\frac{1}{2\pi}\int \hat{f}(\theta)e^{-ik\theta}d\theta.$$
By the inversion formula, we have $$f(\theta)=\sum_{k\in\mathbb{Z}}\hat{f}(k)e^{ik\theta}.$$
Now take $f(\theta)=|\theta|$, $\theta\in [-\pi,\pi]$. Note that $f\in Lip(S^{1})$
We have $$ \hat{f}(k) = \left\{ \begin{array}{rl} \frac{\pi}{2} &\mbox{ if $k=0$} \\ 0 &\mbox{ if $|k|\neq 0$ and $|k|$ is even} \\ -\frac{2}{k^{2}\pi} &\mbox{if $|k|\neq 0$ and $|k|$ is odd} \end{array} \right. $$
Using the inversion formula, we have on $\theta=0$ that $$0=\sum_{k\in\mathbb{Z}}\hat{f}(k).$$
Then,
\begin{eqnarray} 0 &=& \frac{\pi}{2}-\sum_{k\in\mathbb{Z}\ |k|\ odd}\frac{2}{k^{2}\pi} \nonumber \\ &=& \frac{\pi}{2}-\sum_{k\in\mathbb{N}\ |k|\ odd}\frac{4}{k^{2}\pi} \nonumber \\ \end{eqnarray}
This implies $$\sum_{k\in\mathbb{N}\ |k|\ odd}\frac{1}{k^{2}} =\frac{\pi^{2}}{8}$$
If we multiply the last equation by $\frac{1}{2^{2n}}$ with $n=0,1,2,...$ ,we get $$\sum_{k\in\mathbb{N}\ |k|\ odd}\frac{1}{(2^{n}k)^{2}} =\frac{\pi^{2}}{2^{2n}8}$$
Now $$\sum_{n=0,1,...}(\sum_{k\in\mathbb{N}\ |k|\ odd}\frac{1}{(2^{n}k)^{2}}) =\sum_{n=0,1,...}\frac{\pi^{2}}{2^{2n}8}$$
The sum in the left is equal to: $\sum_{k\in\mathbb{N}}\frac{1}{k^{2}}$
The sum in the right is equal to :$\frac{\pi^{2}}{6}$
So we conclude: $$\sum_{k\in\mathbb{N}}\frac{1}{k^{2}}=\frac{\pi^{2}}{6}$$
Note: This is problem 9, Page 208 from the boof of Michael Eugene Taylor - Partial Differential Equation Volume 1.

user91500
Theorem: Let $\lbrace a_n\rbrace$ be a nonincreasing sequence of positive numbers such that $\sum a_n^2$ converges. Then both series $$s:=\sum_{n=0}^\infty(-1)^na_n,\,\delta_k:=\sum_{n=0}^\infty a_na_{n+k},\,k\in\mathbb N $$ converge. Morevere $\Delta:=\sum_{k=1}^\infty(-1)^{k-1}\delta_k$ also converges, and we have the formula $$\sum_{n=0}^\infty a_n^2=s^2+2\Delta.$$ Proof: Knopp. Konrad, Theory and Application of Infinite Series, page 323.
If we let $a_n=\frac1{2n+1}$ in this theorem, then we have $$s=\sum_{n=0}^\infty(-1)^n\frac1{2n+1}=\frac\pi 4$$ $$\delta_k=\sum_{n=0}^\infty\frac1{(2n+1)(2n+2k+1)}=\frac1{2k}\sum_{n=0}^\infty\left(\frac1{2n+1}-\frac1{2n+2k+1}\right)=\frac{1}{2k}\left(1+\frac1 3+...+\frac1 {2k-1}\right)$$ Hence, $$\sum_{n=0}^\infty\frac1{(2n+1)^2}=\left(\frac\pi 4\right)^2+\sum_{k=1}^\infty\frac{(-1)^{k-1}}{k}\left(1+\frac1 3+...+\frac1 {2k-1}\right)=\frac{\pi^2}{16}+\frac{\pi^2}{16}=\frac{\pi^2}{8}$$ and now $$\zeta(2)=\frac4 3\sum_{n=0}^\infty\frac1{(2n+1)^2}=\frac{\pi^2}6.$$

Markus Scheuer
Here's a proof based upon periods and the fact that $\zeta(2)$ and $\frac{\pi^2}{6}$ are periods forming an accessible identity.
The definition of periods below and the proof is from the fascinating introductory survey paper about periods by M. Kontsevich and D. Zagier.
Periods are defined as complex numbers whose real and imaginary parts are values of absolutely convergent integrals of rational functions with rational coefficient over domains in $\mathbb{R}^n$ given by polynomial inequalities with rational coefficients.
The set of periods is therefore a countable subset of the complex numbers. It contains the algebraic numbers, but also many of famous transcendental constants.
In order to show the equality $\zeta(2)=\frac{\pi^2}{6}$ we have to show that both are periods and that $\zeta(2)$ and $\frac{\pi^2}{6}$ form a so-called accessible identity.
First step of the proof: $\zeta(2)$ and $\pi$ are periods
There are a lot of different proper representations of $\pi$ showing that this constant is a period. In the referred paper above following expressions (besides others) of $\pi$ are stated:
\begin{align*} \pi= \iint \limits_{x^2+y^2\leq 1}dxdy=\int_{-\infty}^{\infty}\frac{dx}{1+x^2} \end{align*}
showing that $\pi$ is a period. The known representation
\begin{align*} \zeta(2)=\iint_{0<x<y<1} \frac{dxdy}{(1-x)y} \end{align*}
shows that $\zeta(2)$ is also a period.
$$ $$
Second step: $\zeta(2)$ and $\frac{\pi^2}{6}$ form an accessible identity.
An accessible identity between two periods $A$ and $B$ is given, if we can transform the integral representation of period $A$ by application of the three rules: Additivity (integrand and domain), Change of variables and Newton-Leibniz formula to the integral representation of period $B$.
This implies equality of the periods and the job is done.
In order to show that $\zeta(2)$ and $\frac{\pi^2}{6}$ are accessible identities we start with the integral $I$
$$I=\int_{0}^{1}\int_{0}^{1}\frac{1}{1-xy}\frac{dxdy}{\sqrt{xy}}$$
Expanding $1/(1-xy)$ as a geometric series and integrating term-by-term,
we find that
$$I=\sum_{n=0}^{\infty}\left(n+\frac{1}{2}\right)^{-2}=(4-1)\zeta(2),$$
providing another period representation of $\zeta(2)$.
Changing variables:
$$x=\xi^2\frac{1+\eta^2}{1+\xi^2},\qquad\qquad y=\eta^2\frac{1+\xi^2}{1+\eta^2}$$
with Jacobian $\left|\frac{\partial(x,y)}{\partial(\xi,\eta)}\right|=\frac{4\xi\eta(1-\xi^2\eta^2)}{(1+\xi^2)(1+\eta^2)} =4\frac{(1-xy)\sqrt{xy}}{(1+\xi^2)(1+\eta^2)}$, we find
$$I=4\iint_{0<\eta,\xi\leq 1}\frac{d\xi}{1+\xi^2}\frac{d\eta}{1+\eta^2} =2\int_{0}^{\infty}\frac{d\xi}{1+\xi^2}\int_{0}^{\infty}\frac{d\eta}{1+\eta^2},$$
the last equality being obtained by considering the involution $(\xi,\eta) \mapsto (\xi^{-1},\eta^{-1})$ and comparing this with the last integral representation of $\pi$ above we obtain: $$I=\frac{\pi^2}{2}$$
So, we have shown that $\frac{\pi^2}{6}$ and $\zeta(2)$ are accessible identities and equality follows.

🐢

I Want To Remain Anonymous
As taken from my upcoming textbook:
There is yet another solution to the Basel problem as proposed by Ritelli (2013). His approach is similar to the one by Apostol (1983), where he arrives at
$$\sum_{n\geq1}\frac{1}{n^2}=\frac{\pi^2}{6}\tag1$$
by evaluating the double integral
$$\int_0^1\int_0^1\dfrac{\mathrm{d}x\,\mathrm{d}y}{1-xy}.\tag2$$
Ritelli evaluates in this case the definite integral shown in $(4)$. The starting point comes from realizing that $(1)$ is equivalent to
$$\sum_{n\geq0}\frac{1}{(2n+1)^2}=\frac{\pi^2}{8}\tag3$$
To evaluate the above sum we consider the definite integral
$$\int_0^\infty\int_0^\infty\frac{\mathrm{d}x\,\mathrm{d}y}{(1+y)(1+x^2y)}.\tag4$$
We evaluate $(4)$ first with respect to $x$ and then to $y$
$$\begin{align} \int_0^\infty\left(\frac{1}{1+y}\int_0^\infty\frac{\mathrm{d}x}{1+x^2y}\right)\mathrm{d}y &=\int_0^\infty\left(\frac{1}{1+y}\left[\frac{\tan^{-1}(\sqrt{y}\,x)}{\sqrt{y}}\right]_{x=0}^{x=\infty}\right)\mathrm{d}y\\ &=\frac\pi2\int_0^\infty\frac{\mathrm{d}y}{\sqrt{y}(1+y)}\\ &=\frac\pi2\int_0^\infty\frac{2u}{u(1+u^2)}\mathrm{d}u=\frac{\pi^2}{2},\tag5 \end{align}$$
where we used the substitution $y\leadsto u^2$ in the last step. If we reverse the order of integration one gets
$$\begin{align} \int_0^\infty\left(\int_0^\infty\frac{\mathrm{d}y}{(1+y)(1+x^2y)}\right)\mathrm{d}x&=\int_0^\infty\frac{1}{1-x^2}\left(\int_0^\infty\left(\frac{1}{1+y}-\frac{x^2}{1+x^2y}\right)\mathrm{d}y\right)\mathrm{d}x\\ &=\int_0^\infty\frac{1}{1-x^2}\ln\frac1{x^2}\mathrm{d}x=2\int_0^\infty\frac{\ln x}{x^2-1}\mathrm{d}x.\tag6 \end{align}$$
Hence since $(5)$ and $(6)$ are the same, we have
$$\int_0^\infty\frac{\ln x}{x^2-1}\mathrm{d}x=\frac{\pi^2}{4}.\tag7$$
Furthermore
$$\begin{align} \int_0^\infty\frac{\ln x}{x^2-1}\mathrm{d}x&=\int_0^1\frac{\ln x}{x^2-1}\mathrm{d}x+\int_1^\infty\frac{\ln x}{x^2-1}\mathrm{d}x\\ &=\int_0^1\frac{\ln x}{x^2-1}\mathrm{d}x+\int_0^1\frac{\ln u}{u^2-1}\mathrm{d}u,\tag8 \end{align}$$
where we used the substitution $x\leadsto1/u$. Combining $(7)$ and $(8)$ yields
$$\int_0^1\frac{\ln x}{x^2-1}\mathrm{d}x=\frac{\pi^2}{8}.\tag{9}$$
By expanding the denominator of the integrand in $(10)$ into a geometric series and using the Monotone Convergence Theorem,
$$\int_0^1\frac{\ln x}{x^2-1}\mathrm{d}x=\int_0^1\frac{-\ln x}{1-x^2}\mathrm{d}x=\sum_{n\ge0}\int_0^1(-x^{2n}\ln x)\mathrm{d}x.\tag{10}$$
Using integration by parts one can see that
$$\int_0^1(-x^{2n}\ln x)\mathrm{d}x=\left[-\frac{x^{2n+1}}{2n+1}\ln x\right]^1_0+\int_0^1\frac{x^{2n}}{2n+1}\mathrm{d}x=\frac{1}{(2n+1)^2}\tag{11}$$
Hence from $(10)$, and $(11)$
$$\int_0^1\frac{\ln x}{x^2-1}\mathrm{d}x=\sum_{n\geq0}\frac{1}{(2n+1)^2},\tag{12}$$
which finishes the proof. $$\tag*{$\square$}$$
References:
Daniele Ritelli (2013), Another Proof of $\zeta(2)=\frac{\pi^2}{6}$ Using Double Integrals, The American Mathematical Monthly, Vol. 120, No. 7, pp. 642-645
T. Apostol (1983), A proof that Euler missed: Evaluating $\zeta(2)$ the easy way, Math. Intelligencer 5, pp. 59-60, available at http://dx.doi.org/10.1007/BF03026576.

Eugene Shvarts
This popped up in some reading I'm doing for my research, so I thought I'd contribute! It's a more general twist on the usual pointwise-convergent Fourier series argument.
Consider the eigenvalue problem for the negative Laplacian $\mathcal L$ on $[0,1]$ with Dirichlet boundary conditions; that is, $\mathcal L f:=-f_n'' = \lambda_n f_n$ with $f_n(0) = f_n(1) = 0$. Through inspection we can find that the admissible eigenvalues are $\lambda_n = n^2\pi^2$ for $n=1,2,\ldots$
One can verify that the integral operator $\mathcal Gf(x) = \int_0^1 G(x,y)f(y)\,dy$, where $$G(x,y) = \min(x,y) - xy = \frac{1}{2}\left( -|x-y| + x(1-y) + y(1-x) \right)~~,$$ inverts the negative Laplacian, in the sense that $\mathcal L \mathcal G f = \mathcal G \mathcal L f = f$ on the admissible class of functions (twice weakly differentiable, satisfying the boundary conditions). That is, $G$ is the Green's function for the Dirichlet Laplacian. Because $\mathcal G$ is a self-adjoint, compact operator, we can form an orthonormal basis for $L^2([0,1])$ from its eigenfunctions, and so may express its trace in two ways: $$ \sum_n <f_n,\mathcal G f_n> = \sum_n \frac{1}{\lambda_n} $$and $$\sum_n <f_n,\mathcal G f_n> = \int_0^1 \sum_n f_n(x) <G(x,\cdot),f_n>\,dx = \int_0^1 G(x,x)\,dx~~.$$
The latter quantity is $$ \int_0^1 x(1-x)\,dx = \frac 1 2 - \frac 1 3 = \frac 1 6~~.$$
Hence, we have that $$\sum_n \frac 1 {n^2\pi^2} = \frac 1 6~~\text{, or}~~ \sum_n \frac 1 {n^2} = \frac {\pi^2} 6~~.$$

Markus Scheuer
Here is Euler's Other Proof by Gerald Kimble
> $$
\begin{align*} \frac{\pi^2}{6}&=\frac{4}{3}\frac{(\arcsin 1)^2}{2}\\ &=\frac{4}{3}\int_0^1\frac{\arcsin x}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}\int_0^1\frac{x+\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!}\frac{x^{2n+1}}{2n+1}}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}\int_0^1\frac{x}{\sqrt{1-x^2}}\,dx +\frac{4}{3}\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!(2n+1)}\int_0^1x^{2n}\frac{x}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}+\frac{4}{3}\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!(2n+1)}\left[\frac{(2n)!!}{(2n+1)!!}\right]\\ &=\frac{4}{3}\sum_{n=0}^{\infty}\frac{1}{(2n+1)^2}\\ &=\frac{4}{3}\left(\sum_{n=1}^{\infty}\frac{1}{n^2}-\frac{1}{4}\sum_{n=1}^{\infty}\frac{1}{n^2}\right)\\ &=\sum_{n=1}^{\infty}\frac{1}{n^2} \end{align*}
$$

B_Scheiner
Consider the function $\pi \cot(\pi z)$ which has poles at $z=\pm n$ where n is an integer. Using the L'hopital rule you can see that the residue at these poles is 1.
Now consider the integral $\int_{\gamma_N} \frac{\pi\cot(\pi z)}{z^2} dz$ where the contour $\gamma_N$ is the rectangle with corners given by ±(N + 1/2) ± i(N + 1/2) so that the contour avoids the poles of $\cot(\pi z)$. The integral is bouond in the following way: $\int_{\gamma_N} |\frac{\pi\cot(\pi z)}{z^2} |dz\le Max |(\frac{\pi\cot(\pi z)}{z^2}) | Length(\gamma_N)$. It can easily be shown that on the contour $\gamma_N$ that $\pi \cot(\pi z)< M$ where M is some constant. Then we have
$\int_{\gamma_N} |\frac{\pi\cot(\pi z)}{z^2} |dz\le M Max |\frac{1}{z^2} | Length(\gamma_N) = (8N+4) \frac{M}{\sqrt{2(1/2+N)^2}^2}$
where (8N+4) is the lenght of the contour and $\sqrt{2(1/2+N)^2}$ is half the diagonal of $\gamma_N$. In the limit that N goes to infinity the integral is bound by 0 so we have $\int_{\gamma_N} \frac{\pi\cot(\pi z)}{z^2} dz =0$
by the cauchy residue theorem we have 2πiRes(z = 0) + 2πi$\sum$Residues(z$\ne$ 0) = 0. At z=0 we have Res(z=0)=$-\frac{\pi^2}{3}$, and $Res (z=n)=\frac{1}{n^2}$ so we have
$2\pi iRes(z = 0) + 2\pi i\sum Residues(z\ne 0) = -\frac{\pi^2}{3}+2\sum_{1}^{\infty} \frac{1}{n^2} =0$
Where the 2 in front of the residue at n is because they occur twice at +/- n.
We now have the desired result $\sum_{1}^{\infty} \frac{1}{n^2}=\frac{\pi^2}{6}$.

Meadara
I saw this proof in an extract of the College Mathematics Journal.
Consider the Integeral : $I = \int_0^{\pi/2}\ln(2\cos x)dx$
From $2\cos(x) = e^{ix} + e^{-ix}$ , we have:
$$\int_0^{\pi/2}\ln\left(e^{ix} + e^{-ix}\right)dx = \int_0^{\pi/2}\ln\left(e^{ix}(1 + e^{-2ix})\right)dx=\int_0^{\pi/2}ixdx + \int_0^{\pi/2}\ln(1 + e^{-2ix})dx$$
The Taylor series expansion of $\ln(1+x)=x -\frac{x^2}{2} +\frac{x^3}{3}-\cdots$
Thus , $\ln(1+e^{-2ix}) = e^{-2ix}- \frac{e^{-4ix}}{2} + \frac{e^{-6ix}}{3} - \cdots $, then for $I$ :
$$I = \frac{i\pi^2}{8}+\left[-\frac{e^{-2ix}}{2i}+\frac{e^{-4ix}}{2\cdot 4i}-\frac{e^{-6ix}}{3\cdot 6i}-\cdots\right]_0^\frac{\pi}{2}$$
$$I = \frac{i\pi^2}{8}-\frac{1}{2i}\left[\frac{e^{-2ix}}{1^2}-\frac{e^{-4ix}}{2^2}+\frac{e^{-6ix}}{3^2}-\cdots\right]_0^\frac{\pi}{2}$$
By evaluating we get something like this..
$$I = \frac{i\pi^2}{8}-\frac{1}{2i}\left[\frac{-2}{1^2}-\frac{0}{2^2}+\frac{-2}{3^2}-\cdots\right]_0^\frac{\pi}{2}$$
Hence
$$\int_0^{\pi/2}\ln(2\cos x)dx=\frac{i\pi^2}{8}-i\sum_{k=0}^\infty \frac{1}{(2k+1)^2}$$
So now we have a real integral equal to an imaginary number, thus the value of the integral should be zero.
Thus, $\sum_{k=0}^\infty \frac{1}{(2k+1)^2}=\frac{\pi^2}{8}$
But let $\sum_{k=0}^\infty \frac{1}{k^2}=E$ .We get $\sum_{k=0}^\infty \frac{1}{(2k+1)^2}=\frac{3}{4} E$
And as a result $$\sum_{k=0}^\infty \frac{1}{k^2} = \frac{\pi^2}{6}$$

dustin
I have another method as well. From skimming the previous solutions, I don't think it is a duplicate of any of them
In Complex analysis, we learn that $\sin(\pi z) = \pi z\Pi_{n=1}^{\infty}\Big(1 - \frac{z^2}{n^2}\Big)$ which is an entire function with simple zer0s at the integers. We can differentiate term wise by uniform convergence. So by logarithmic differentiation we obtain a series for $\pi\cot(\pi z)$. $$ \frac{d}{dz}\ln(\sin(\pi z)) = \pi\cot(\pi z) = \frac{1}{z} - 2z\sum_{n=1}^{\infty}\frac{1}{n^2 - z^2} $$ Therefore, $$ -\sum_{n=1}^{\infty}\frac{1}{n^2 - z^2} = \frac{\pi\cot(\pi z) - \frac{1}{z}}{2z} $$ We can expand $\pi\cot(\pi z)$ as $$ \pi\cot(\pi z) = \frac{1}{z} - \frac{\pi^2}{3}z - \frac{\pi^4}{45}z^3 - \cdots $$ Thus, \begin{align} \frac{\pi\cot(\pi z) - \frac{1}{z}}{2z} &= \frac{- \frac{\pi^2}{3}z - \frac{\pi^4}{45}z^3-\cdots}{2z}\\ -\sum_{n=1}^{\infty}\frac{1}{n^2 - z^2}&= -\frac{\pi^2}{6} - \frac{\pi^4}{90}z^2 - \cdots\\ -\lim_{z\to 0}\sum_{n=1}^{\infty}\frac{1}{n^2 - z^2}&= \lim_{z\to 0}\Big(-\frac{\pi^2}{6} - \frac{\pi^4}{90}z^2 - \cdots\Big)\\ -\sum_{n=1}^{\infty}\frac{1}{n^2}&= -\frac{\pi^2}{6}\\ \sum_{n=1}^{\infty}\frac{1}{n^2}&= \frac{\pi^2}{6} \end{align}

Elias
See evaluations of Riemann Zeta Function $\zeta(2)=\sum_{n=1}^\infty\frac{1}{n^2}$ in mathworld.wolfram.com and a solution by in D. P. Giesy in Mathematics Magazine:
D. P. Giesy, Still another elementary proof that $\sum_{n=1}^\infty \frac{1}{n^2}=\frac{\pi^2}{6}$, Math. Mag. 45 (1972) 148-149.
Unfortunately I did not get a link to this article. But there is a link to a note from Robin Chapman seems to me a variation of proof's Giesy.

barto
Applying the usual trick 1 transforming a series to an integral, we obtain
$$\sum_{n=1}^\infty\frac1{n^2}=\int_0^1\int_0^1\frac{dxdy}{1-xy}$$
where we use the Monotone Convergence Theorem to integrate term-wise.
Then there's this ingenious change of variables 2, which I learned from Don Zagier during a lecture, and which he in turn got from a colleague:
$$(x,y)=\left(\frac{\cos v}{\cos u},\frac{\sin u}{\sin v}\right),\quad0\leq u\leq v\leq \frac\pi2$$
One verifies that it is bijective between the rectangle $[0,1]^2$ and the triangle $0\leq u\leq v\leq \frac\pi2$, and that its Jacobian determinant is precisely $1-x^2y^2$, which means $\frac1{1-x^2y^2}$ would be a neater integrand. For the moment, we have found
$$J=\int_0^1\int_0^1\frac{dxdy}{1-x^2y^2}=\frac{\pi^2}8$$ (the area of the triangular domain in the $(u,v)$ plane).
There are two ways to transform $\int\frac1{1-xy}$ into something $\int\frac1{1-x^2y^2}$ish:
Manipulate $S=\sum_{n=1}^\infty\frac1{n^2}$: We have $\sum_{n=1}^\infty\frac1{(2n)^2}=\frac14S$ so $\sum_{n=0}^\infty\frac1{(2n+1)^2}=\frac34S$. Applying the series-integral transformation, we get $\frac34S=J$ so $$S=\frac{\pi^2}6$$
Manipulate $I=\int_0^1\int_0^1\frac{dxdy}{1-xy}$: Substituting $(x,y)\leftarrow(x^2,y^2)$ we have $I=\int_0^1\int_0^1\frac{4xydxdy}{1-x^2y^2}$ so $$J=\int_0^1\int_0^1\frac{dxdy}{1-x^2y^2}=\int_0^1\int_0^1\frac{(1+xy-xy)dxdy}{1-x^2y^2}=I-\frac14I$$ whence $$I=\frac43J=\frac{\pi^2}6$$
(It may be seen that they are essentially the same methods.)
After looking at the comments it seems that this looks a lot like Proof 2 in the article by R. Chapman.
See also: Multiple Integral $\int\limits_0^1\!\!\int\limits_0^1\!\!\int\limits_0^1\!\!\int\limits_0^1\frac1{1-xyuv}\,dx\,dy\,du\,dv$
1 See e.g. Proof 1 in Chapman's article.
2 It may have been a different one; maybe as in the above article. Either way, the idea to do something trigonometric was not mine.

Asier Calbet
The sum can be written as the integral: $$\int_0^{\infty} \frac{x}{e^x-1} dx $$ This integral can be evaluated using a rectangular contour from 0 to $\infty$ to $\infty + \pi i$ to $ 0$ .

John M. Campbell
There is a simple way of proving that $\sum_{n=1}^{\infty}\frac{1}{n^2} = \frac{\pi^2}{6}$ using the following well-known series identity: $$\left(\sin^{-1}(x)\right)^{2} = \frac{1}{2}\sum_{n=1}^{\infty}\frac{(2x)^{2n}}{n^2 \binom{2n}{n}}.$$ From the above equality, we have that $$x^2 = \frac{1}{2}\sum_{n=1}^{\infty}\frac{(2 \sin(x))^{2n}}{n^2 \binom{2n}{n}},$$ and we thus have that: $$\int_{0}^{\pi} x^2 dx = \frac{\pi^3}{12} = \frac{1}{2}\sum_{n=1}^{\infty}\frac{\int_{0}^{\pi} (2 \sin(x))^{2n} dx}{n^2 \binom{2n}{n}}.$$ Since $$\int_{0}^{\pi} \left(\sin(x)\right)^{2n} dx = \frac{\sqrt{\pi} \ \Gamma\left(n + \frac{1}{2}\right)}{\Gamma(n+1)},$$ we thus have that: $$\frac{\pi^3}{12} = \frac{1}{2}\sum_{n=1}^{\infty}\frac{ 4^{n} \frac{\sqrt{\pi} \ \Gamma\left(n + \frac{1}{2}\right)}{\Gamma(n+1)} }{n^2 \binom{2n}{n}}.$$ Simplifying the summand, we have that $$\frac{\pi^3}{12} = \frac{1}{2}\sum_{n=1}^{\infty}\frac{\pi}{n^2},$$ and we thus have that $\sum_{n=1}^{\infty}\frac{1}{n^2} = \frac{\pi^2}{6}$ as desired.
`
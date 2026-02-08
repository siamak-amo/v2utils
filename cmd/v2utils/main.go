/* V2utils provides xray-core compatible utilities
   Copyright 2025-2026 Ahmad <edu.siamak@gmail.com>

   V2utils is free software: you can redistribute it and/or modify it
   under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   V2utils is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

// V2utils current version
const Version = "1.4";


func main() {
	// Initializing user options
	// It calls exit() on fatal failures, so
	// we don't need to handle failures.
	opt := init_opt ();

	// Running the main loop
	main_loop (opt);
}

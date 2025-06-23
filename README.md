# wsinging-box

The universal singing platform 🎶.

This is a customized version of sing-box for my personal usage, it might contain undocumented features (or bugs) which do not exist in the official version.

I will not hold any responsibility in instability of this version.

Build command used:
```ps1
go build -trimpath -o dist/sing-box.exe -ldflags '-s -buildid= -X github.com/sagernet/sing-box/constant.Version=1.0.0' ./cmd/sing-box
```

## Documentation

https://sing-box.sagernet.org

## License

```
Copyright (C) 2022 by nekohasekai <contact-sagernet@sekai.icu>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

In addition, no derivative work may use the name or imply association
with this application without prior consent.
```
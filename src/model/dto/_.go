package dto

//////////////////////////////////////////////////////////////////////////////////////
// OtherDto

// @Model         _ImageDto
// @Description   上传图片信息
// @Property      url  string  true "图片链接"
// @Property      size integer true "图片大小，单位为字节"

//////////////////////////////////////////////////////////////////////////////////////
// Result

// @Model         Result<UserDto>
// @Description   返回用户信息
// @Property      code    integer           true "返回响应码"
// @Property      message string            true "返回信息"
// @Property      data    object(#_UserDto) true "返回数据"

// @Model         Result<LoginDto>
// @Description   登录信息
// @Property      code    integer            true "返回响应码"
// @Property      message string             true "返回信息"
// @Property      data    object(#_LoginDto) true "返回数据"

// @Model         Result<UserExtraDto>
// @Description   返回用户与数量信息
// @Property      code     integer                   true "返回响应码"
// @Property      message  string                    true "返回信息"
// @Property      data     object(#_UserAndExtraDto) true "返回数据"

// @Model         Result<VideoDto>
// @Description   返回视频信息
// @Property      code    integer            true "返回响应码"
// @Property      message string             true "返回信息"
// @Property      data    object(#_VideoDto) true "返回数据"

// @Model         Result<ImageDto>
// @Description   返回上传图片信息
// @Property      code    integer            true "返回响应码"
// @Property      message string             true "返回信息"
// @Property      data    object(#_ImageDto) true "返回数据"

//////////////////////////////////////////////////////////////////////////////////////
// _Page

// @Model         Page<UserDto>
// @Description   用户分页信息
// @Property      total  integer          true "数据总量"
// @Property      page   int              true "当前页"
// @Property      limit  int              true "页大小"
// @Property      data   array(#_UserDto) true "返回数据"

// @Model         Page<VideoDto>
// @Description   视频分页信息
// @Property      total  integer           true "数据总量"
// @Property      page   int               true "当前页"
// @Property      limit  int               true "页大小"
// @Property      data   array(#_VideoDto) true "返回数据"

//////////////////////////////////////////////////////////////////////////////////////
// PageResult

// @Model         Result<Page<UserDto>>
// @Description   返回用户分页信息
// @Property      code      integer                true "返回响应码"
// @Property      message   string                 true "返回信息"
// @Property      data      object(#Page<UserDto>) true "返回数据"

// @Model         Result<Page<VideoDto>>
// @Description   返回视频分页信息
// @Property      code      integer                 true "返回响应码"
// @Property      message   string                  true "返回信息"
// @Property      data      object(#Page<VideoDto>) true "返回数据"

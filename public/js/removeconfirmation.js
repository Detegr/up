$(function()
{
	$(document.body).click(function() {
		$(".removeConfirmButton:visible").slideToggle("fast", function() {
			$(".removeButton:hidden").slideToggle("fast");
		});
	});
	$(".removeButton").click(function() {
		var $this=$(this);
		$this.slideToggle("fast", function() {
			$this.parent().find(".removeConfirmButton").slideToggle("fast");
		});
		return false;
	});
});
